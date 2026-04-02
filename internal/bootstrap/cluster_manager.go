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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	rpc     *grpc.ClientConn
	running bool

	subscriptions map[string]core.TopicListener

	cSub     chan Sub
	cInbound chan *protocol.Topic
}

func (c *ClusterManager) HashRing(r core.RingRequest) {
	if !c.running {
		r.Async <- []core.Node{}
		return
	}
	dsp := protocol.NewPostofficeServiceClient(c.rpc)
	stream, err := dsp.HashRing(context.Background(), &protocol.Request{Prefix: 0})
	if err != nil {
		return
	}
	ring := make([]core.Node, 0)
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			core.AppLog.Debug().Msgf("streaming error %s", err.Error())
			break
		}
		ring = append(ring, core.Node{Name: data.Name, RingToken: data.Hash, RpcEndpoint: data.Endpoint, IP: data.Address})
	}
	r.Async <- ring
}

func (c *ClusterManager) KeyRing(r core.RingRequest) {
	if !c.running {
		r.Async <- []core.Node{}
		return
	}
	dsp := protocol.NewPostofficeServiceClient(c.rpc)
	stream, err := dsp.KeyRing(context.Background(), &protocol.Request{Prefix: r.Token})
	if err != nil {
		return
	}
	ring := make([]core.Node, 0)
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			core.AppLog.Debug().Msgf("streaming error %s", err.Error())
			break
		}
		ring = append(ring, core.Node{Name: data.Name, RingToken: data.Hash, RpcEndpoint: data.Endpoint, IP: data.Address})
	}
	r.Async <- ring
}

func (c *ClusterManager) RingToken(key []byte) uint32 {
	return util.Hash(key)
}

func (c *ClusterManager) Request(r core.DataRequest) {
	if !c.running {
		r.Async <- core.Chunk{Remaining: false}
		return
	}
	dsp := protocol.NewPostofficeServiceClient(c.rpc)
	req := protocol.Request{Prefix: r.Prefix, Opt: r.Opt, Data: &protocol.Data{Key: r.Key, Value: r.Value, Header: &protocol.Header{Revision: r.Revision, FactoryId: r.FactoryId, ClassId: r.ClassId, Mutable: r.Mutable}}}
	if r.Opt == core.QUERY_DATA_REQUEST || r.Opt == core.PULL_DATA_REQUEST {
		dt, err := event.Export(r.Criteria, core.COMPOSIT_KEY_MAX)
		if err != nil {
			r.Async <- core.Chunk{Remaining: false, Data: protocol.Response{Successful: false, Message: err.Error()}}
			return
		}
		q := protocol.Query{Id: r.Criteria.QTopic(), Criteria: dt}
		req.Query = &q
	}
	stream, err := dsp.Request(context.Background(), &req)
	if err != nil {
		r.Async <- core.Chunk{Remaining: false, Data: protocol.Response{Successful: false, Message: err.Error()}}
		return
	}
	crt := core.Chunk{Remaining: false}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			core.AppLog.Warn().Msgf("streaming error %s", err.Error())
			crt.Data = err
			break
		}
		if resp.Successful {
			r.Async <- core.Chunk{Remaining: true, Data: resp}
		}
	}
	r.Async <- crt
}

func (c *ClusterManager) Publish(e *protocol.Topic) error {
	if !c.running {
		return fmt.Errorf("not started")
	}
	dsp := protocol.NewPostofficeServiceClient(c.rpc)
	resp, err := dsp.Publish(context.Background(), e)
	if err != nil {
		return err
	}
	core.AppLog.Debug().Msgf("topic publish %v", resp)
	return nil
}

//func (c *ClusterManager) List(query core.Query) {
//req := core.DataRequest{Opt: core.QUERY_DATA_REQUEST, Criteria: query}
//req.Async = query.QCc()
//c.Request(req)
//}

func (c *ClusterManager) Subscribe(topic string, listener core.TopicListener) error {
	if !c.running {
		return fmt.Errorf("not started")
	}
	dsp := protocol.NewPostofficeServiceClient(c.rpc)
	resp, err := dsp.Subscribe(context.Background(), &protocol.Topic{NodeId: c.App.NodeId(), Tag: c.App.Context(), Name: topic})
	if err != nil {
		return err
	}
	c.cSub <- Sub{opt: OPT_SUB, name: topic, listener: listener}
	core.AppLog.Debug().Msgf("topic registered %v", resp)
	return nil
}

func (c *ClusterManager) Unsubscribe(topic string) error {
	if !c.running {
		return fmt.Errorf("not started")
	}
	dsp := protocol.NewPostofficeServiceClient(c.rpc)
	resp, err := dsp.Unsubscribe(context.Background(), &protocol.Topic{Tag: c.App.Context(), Name: topic})
	if err != nil {
		return err
	}
	c.cSub <- Sub{opt: OPT_UNSUB, name: topic}
	core.AppLog.Debug().Msgf("topic unregistered %v", resp)
	return nil
}

func (c *ClusterManager) disconnect() error {
	if !c.running {
		return fmt.Errorf("not started")
	}
	c.running = false
	dsp := protocol.NewPostofficeServiceClient(c.rpc)
	resp, err := dsp.Disconnect(context.Background(), &protocol.Topic{Tag: c.App.Context(), NodeId: c.App.NodeId()})
	if err != nil {
		return err
	}
	core.AppLog.Debug().Msgf("disconnecting topic %v", resp)
	return c.rpc.Close()
}

func (c *ClusterManager) connect(host string) error {
	retries := RPC_CONNECT_RETRIES
	for {
		tcp, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, core.RPC_PORT), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			retries--
			if retries > 0 {
				core.AppLog.Warn().Msgf("retrying to connect gprc %s with retried times %d", err.Error(), retries)
				time.Sleep(3 * time.Second)
				continue
			}
			return err
		}
		c.rpc = tcp
		break
	}
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
ro:
	dsp := protocol.NewPostofficeServiceClient(c.rpc)
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
