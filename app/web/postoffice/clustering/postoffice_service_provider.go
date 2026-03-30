package clustering

import (
	context "context"
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
)

func (c *DataServiceProvider) HashRing(request *protocol.Request, stream grpc.ServerStreamingServer[protocol.HashNode]) error {
	rq := make(chan []core.Node, 1)
	defer close(rq)
	c.Mll.rangeRing(core.RingRequest{Async: rq, Opt: ALL_RING_OPT})
	ring := <-rq
	for _, n := range ring {
		hn := protocol.HashNode{Hash: n.RingToken, Endpoint: n.RpcEndpoint, Name: n.Name, Address: n.IP}
		if err := stream.Send(&hn); err != nil {
			return err
		}
	}
	return nil
}

func (c *DataServiceProvider) KeyRing(request *protocol.Request, stream grpc.ServerStreamingServer[protocol.HashNode]) error {
	rq := make(chan []core.Node, 1)
	defer close(rq)
	c.Mll.rangeRing(core.RingRequest{Async: rq, Opt: REPLICA_RING_OPT, Token: request.Prefix})
	ring := <-rq
	for _, n := range ring {
		hn := protocol.HashNode{Hash: n.RingToken, Endpoint: n.RpcEndpoint, Name: n.Name, Address: n.IP}
		if err := stream.Send(&hn); err != nil {
			return err
		}
	}
	return nil
}

func (c *DataServiceProvider) Receive(topic *protocol.Topic, stream grpc.ServerStreamingServer[protocol.Topic]) error {
	aq := make(chan ReceiverAsync, 2)
	c.DRequest <- TopicRequest{Opt: RECEIVER_START, Name: topic.NodeId, Async: aq}
	ch := <-aq
	close(aq)
	core.AppLog.Debug().Msgf("start event receiver on [%s]", topic.NodeId)
	defer close(ch.Rev)
	defer close(ch.Q)
	receiving := true
	for receiving {
		select {
		case <-ch.Q:
			receiving = false
		case resp := <-ch.Rev:
			err := stream.Send(resp)
			if err != nil {
				core.AppLog.Debug().Msgf("send error %s", err.Error())
				receiving = false
			}
		}
	}
	c.DRequest <- TopicRequest{Opt: RECEIVER_REMOVE, Name: topic.NodeId}
	core.AppLog.Debug().Msgf("stop evnt receiver from on [%s]", topic.NodeId)
	c.Mll.MRequest <- core.RingRequest{Opt: SYNC_SUB_OPT, Source: core.RingSync{Sub: core.Subscription{NodeId: topic.NodeId, Deleting: true}}}
	return nil
}
func (c *DataServiceProvider) Disconnect(ctx context.Context, topic *protocol.Topic) (*protocol.Response, error) {
	core.AppLog.Debug().Msgf("receiver disconnected %s", topic.NodeId)
	c.DRequest <- TopicRequest{Opt: RECEIVER_END, Name: topic.NodeId}
	return &protocol.Response{Successful: true, Message: "disconnected"}, nil
}
func (c *DataServiceProvider) Publish(ctx context.Context, in *protocol.Topic) (*protocol.Response, error) {
	return c.runPublish(in)
}

func (c *DataServiceProvider) Subscribe(ctx context.Context, in *protocol.Topic) (*protocol.Response, error) {
	c.Mll.MRequest <- core.RingRequest{Opt: SYNC_SUB_OPT, Source: core.RingSync{Sub: core.Subscription{NodeId: in.NodeId, Tag: in.Tag, Topic: in.Name, Endpoint: c.rpcEndpoint}}}
	return &protocol.Response{Successful: true, Message: "topic created"}, nil
}
func (c *DataServiceProvider) Unsubscribe(ctx context.Context, in *protocol.Topic) (*protocol.Response, error) {
	c.Mll.MRequest <- core.RingRequest{Opt: SYNC_SUB_OPT, Source: core.RingSync{Sub: core.Subscription{NodeId: in.NodeId, Tag: in.Tag, Topic: in.Name, Endpoint: c.rpcEndpoint, Deleting: true}}}
	return &protocol.Response{Successful: true, Message: "topic removed"}, nil
}

func (c *DataServiceProvider) Request(request *protocol.Request, stream grpc.ServerStreamingServer[protocol.Response]) error {
	switch request.Opt {
	case core.CREATE_DATA_REQUEST:
		resp, _ := c.runCreate(request)
		stream.Send(resp)
	case core.GET_DATA_REQUEST:
		rc := make(chan *protocol.Response, 3)
		defer close(rc)
		go c.runGet(request, rc)
		for resp := range rc {
			stream.Send(resp)
			if !resp.Successful {
				break
			}
		}
	case core.QUERY_DATA_REQUEST:
		rc := make(chan *protocol.Response, 3)
		defer close(rc)
		go c.runGet(request, rc)
		for resp := range rc {
			if !resp.Successful {
				break
			}
			stream.Send(resp)
		}
	case core.UPDATE_DATA_REQUEST:
		rc := make(chan *protocol.Response, 3)
		defer close(rc)
		c.runUpdate(request, rc)
		resp := <-rc
		stream.Send(resp)
	case core.DELETE_DATA_REQUEST:
		rc := make(chan *protocol.Response, 3)
		defer close(rc)
		c.runDelete(request, rc)
		resp := <-rc
		stream.Send(resp)
	case core.RESET_DATA_REQUEST:
		rc := make(chan *protocol.Response, 3)
		defer close(rc)
		c.runReset(request, rc)
		resp := <-rc
		stream.Send(resp)
	case core.PULL_DATA_REQUEST:
		rc := make(chan *protocol.Response, 3)
		defer close(rc)
		core.AppLog.Debug().Msgf("run local pull %v", request)
		go c.pull(request.Prefix, request.Prefix, rc)
		for resp := range rc {
			if !resp.Successful {
				break
			}
			stream.Send(resp)
		}
	default:
		stream.Send(&protocol.Response{Successful: false, Message: fmt.Sprintf("request opt not suuported %d", request.Opt)})
	}
	return nil
}
