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
)

type ClusterManager struct {
	App     *AppManager
	rpc     *grpc.ClientConn
	running bool

	subscriptions map[string]core.TopicListener

	CRequest chan core.TopicListener
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
		dt, err := event.Export(r.Criteria, 100)
		if err != nil {
			r.Async <- core.Chunk{Remaining: false, Data: protocol.Response{Successful: false, Message: err.Error()}}
			return
		}
		q := protocol.Query{Id: r.Criteria.QId(), Criteria: dt}
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
		r.Async <- core.Chunk{Remaining: true, Data: resp}
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
func (c *ClusterManager) List(query core.Query) {
	req := core.DataRequest{Opt: core.QUERY_DATA_REQUEST, Criteria: query}
	req.Async = query.QCc()
	c.Request(req)
}

func (c *ClusterManager) Subscribe(topic string, listener core.TopicListener) error {
	if !c.running {
		return fmt.Errorf("not started")
	}
	dsp := protocol.NewPostofficeServiceClient(c.rpc)
	resp, err := dsp.Subscribe(context.Background(), &protocol.Topic{NodeId: c.App.NodeId(), Tag: c.App.Context(), Name: topic})
	if err != nil {
		return err
	}
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

func (s *ClusterManager) connect(host string) error {
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
		s.rpc = tcp
		break
	}
	s.running = true
	go s.receive()
	return nil
}

func (c *ClusterManager) receive() {
	dsp := protocol.NewPostofficeServiceClient(c.rpc)
	stream, err := dsp.Receive(context.Background(), &protocol.Topic{NodeId: c.App.NodeId(), Tag: c.App.Context()})
	if err != nil {
		core.AppLog.Warn().Msgf("rpc connection error %s", err.Error())
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
		core.AppLog.Debug().Msgf("topic %v", resp)
		//c.tpl.OnTopic(resp)
	}
	core.AppLog.Warn().Msg("rpc closed from remote")
}
