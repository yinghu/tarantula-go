package bootstrap

import (
	"context"
	"fmt"
	"io"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	RPC_CONNECT_RETRIES int = 3

	SUB_CHAN_SIZE   int = 3
	TOPIC_CHAN_SIZE int = 12

	OPT_SUB   int = 100
	OPT_UNSUB int = 200
)

type Sub struct {
	opt      int
	name     string
	listener core.TopicListener
}

type ClusterManager struct {
	App     *AppManager
	running bool

	subscriptions map[string]core.TopicListener

	cSub     chan Sub
	cInbound chan *protocol.Topic
	cHost    string
	cPool    core.RpcConnPool
}

func (c *ClusterManager) HashRing(r core.RingRequest) (grpc.ServerStreamingClient[protocol.HashNode], error) {
	if !c.running {
		return nil, fmt.Errorf("cluster not started")
	}
	conn, err := c.cPool.Conn()
	if err != nil {
		return nil, err
	}
	dsp := protocol.NewPostofficeServiceClient(conn.Conn)
	return dsp.HashRing(context.Background(), &protocol.Request{Prefix: 0})
}

func (c *ClusterManager) KeyRing(r core.RingRequest) (grpc.ServerStreamingClient[protocol.HashNode], error) {
	if !c.running {
		return nil, fmt.Errorf("cluster not started")
	}
	conn, err := c.cPool.Conn()
	if err != nil {
		return nil, err
	}
	dsp := protocol.NewPostofficeServiceClient(conn.Conn)
	return dsp.KeyRing(context.Background(), &protocol.Request{Prefix: r.Token})
}

func (c *ClusterManager) RingToken(key []byte) uint32 {
	return util.Hash(key)
}

func (c *ClusterManager) Request(r core.DataRequest) (grpc.ServerStreamingClient[protocol.Response], error) {
	if !c.running {
		return nil, fmt.Errorf("cluster not started")
	}
	req := protocol.Request{Prefix: r.Prefix, Opt: r.Opt, Data: &protocol.Data{Key: r.Key, Value: r.Value, Header: &protocol.Header{Revision: r.Revision, FactoryId: r.FactoryId, ClassId: r.ClassId, Mutable: r.Mutable}}}
	if r.Opt == core.QUERY_DATA_REQUEST || r.Opt == core.PULL_DATA_REQUEST {
		mf, existed := TopicFactoryRegistry[r.Criteria.QTopic()]
		if !existed {
			return nil, fmt.Errorf("topic factory not existed")
		}
		dt, err := mf().Export(r.Criteria)
		if err != nil {
			return nil, err
		}
		q := protocol.Query{Id: r.Criteria.QTopic(), Criteria: dt}
		req.Query = &q
	}
	conn, err := c.cPool.Conn()
	if err != nil {
		return nil, err
	}
	dsp := protocol.NewPostofficeServiceClient(conn.Conn)
	return dsp.Request(context.Background(), &req)
}

func (c *ClusterManager) Publish(e *protocol.Topic) (*protocol.Response, error) {
	if !c.running {
		return &protocol.Response{Successful: false}, fmt.Errorf("not started")
	}
	conn, err := c.cPool.Conn()
	if err != nil {
		return &protocol.Response{Successful: false}, err
	}
	dsp := protocol.NewPostofficeServiceClient(conn.Conn)
	return dsp.Publish(context.Background(), e)
}

func (c *ClusterManager) Subscribe(topic string, listener core.TopicListener) error {
	if !c.running {
		return fmt.Errorf("not started")
	}
	conn, err := c.cPool.Conn()
	if err != nil {
		return err
	}

	dsp := protocol.NewPostofficeServiceClient(conn.Conn)
	resp, err := dsp.Subscribe(context.Background(), &protocol.Topic{NodeId: c.App.NodeId(), Tag: c.App.Context(), Name: topic})
	if err != nil {
		return err
	}
	c.cSub <- Sub{opt: OPT_SUB, name: topic, listener: listener}
	core.AppLog.Debug().Msgf("topic registered %v %s", resp.Successful, topic)
	return nil
}

func (c *ClusterManager) Unsubscribe(topic string) error {
	if !c.running {
		return fmt.Errorf("not started")
	}
	conn, err := c.cPool.Conn()
	if err != nil {
		return err
	}

	dsp := protocol.NewPostofficeServiceClient(conn.Conn)
	resp, err := dsp.Unsubscribe(context.Background(), &protocol.Topic{Tag: c.App.Context(), Name: topic})
	if err != nil {
		return err
	}
	c.cSub <- Sub{opt: OPT_UNSUB, name: topic}
	core.AppLog.Debug().Msgf("topic unregistered %v %s", resp.Successful, topic)
	return nil

}

func (c *ClusterManager) disconnect() error {
	if !c.running {
		return fmt.Errorf("not started")
	}
	c.running = false
	c.cPool.Release()
	return nil
}

func (c *ClusterManager) connect(host string) error {
	c.cHost = fmt.Sprintf("%s:%d", host, core.RPC_PORT)
	c.cPool = core.RpcConnPool{Target: c.cHost, Tag: c.App.Context(), NodeId: c.App.NodeId()}
	c.cPool.Start()
	c.subscriptions = make(map[string]core.TopicListener)
	c.cSub = make(chan Sub, SUB_CHAN_SIZE)
	c.cInbound = make(chan *protocol.Topic, TOPIC_CHAN_SIZE)
	c.running = true
	go c.async()
	go c.receive()
	return nil
}

func (c *ClusterManager) receive() {
	retries := RPC_CONNECT_RETRIES
	conn, err := c.cPool.Conn()
	if err != nil {
		panic(err.Error())
	}
ro:
	dsp := protocol.NewPostofficeServiceClient(conn.Conn)
	stream, err := dsp.Receive(context.Background(), &protocol.Topic{NodeId: c.App.NodeId(), Tag: c.App.Context()})
	if err != nil {
		retries--
		if retries > 0 {
			core.AppLog.Warn().Msgf("rpc connection retry with %s %d", err.Error(), retries)
			time.Sleep(3 * time.Second)
			goto ro
		}
		core.AppLog.Warn().Msgf("rpc connection error after retried %s", err.Error())
		return
	}
	for c.running {
		resp, err := stream.Recv()
		if err == io.EOF {
			core.AppLog.Debug().Msgf("eof %s", err.Error())
			break
		}
		if err != nil {
			core.AppLog.Warn().Msgf("streaming error %s", err.Error())
			break
		}
		c.cInbound <- resp
	}
	core.AppLog.Warn().Msgf("cluster manager receiver closed from remote %v", c.running)
}

func (c *ClusterManager) async() {
	for c.running {
		select {
		case topic := <-c.cInbound:
			tl, ok := c.subscriptions[topic.Name]
			if ok {
				tl.OnTopic(topic)
			} else {
				core.AppLog.Debug().Msgf("dead topic %v", topic)
			}
		case sub := <-c.cSub:
			switch sub.opt {
			case OPT_SUB:
				c.subscriptions[sub.name] = sub.listener
			case OPT_UNSUB:
				delete(c.subscriptions, sub.name)
			}
		}
	}
	core.AppLog.Warn().Msgf("cluster manager async task closed from remote %v", c.running)
	clear(c.subscriptions)
	close(c.cInbound)
	close(c.cSub)
}

func (c *ClusterManager) Forward(level zerolog.Level, log []byte) {
	if !c.running {
		return
	}
	e := protocol.LogEvent{}
	err := protojson.Unmarshal(log, &e)
	if err != nil {
		e.Level = "error"
		e.Message = err.Error()
		e.Time = timestamppb.Now()
		e.Source = "forwarder:325"
	}
	tf := event.LogEventFactory{}
	t, err := tf.FromLogEvent(&e)
	t.NodeId = c.App.NodeId()
	t.Tag = c.App.Context()
	id, _ := c.App.Sequence().Id()
	t.Event.Id = uint64(id)
}
