package clustering

import (
	context "context"
	"fmt"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
)

func (c *DataServiceProvider) HashRing(request *protocol.Request, stream grpc.ServerStreamingServer[protocol.HashNode]) error {
	rq := make(chan []core.Node, 1)
	defer close(rq)
	c.Mll.rangeRing(core.RingRequest{Async: rq, Opt: ALL_RING_OPT})
	ring := <-rq
	for _, n := range ring {
		hn := protocol.HashNode{Hash: n.RingToken, Endpoint: n.RpcEndpoint, Name: n.Name, Address: n.IP, Meta: n.Meta}
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

func (c *DataServiceProvider) Receive(topic *protocol.Topic, stream grpc.ServerStreamingServer[protocol.Mail]) error {
	aq := make(chan ReceiverAsync, 2)
	c.DRequest <- TopicRequest{Opt: RECEIVER_START, Name: topic.NodeId, Async: aq}
	ch := <-aq
	close(aq)
	core.AppLog.Info().Msgf("start event receiver on [%s]", topic.NodeId)
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

func (c *DataServiceProvider) Subscribe(ctx context.Context, in *protocol.Subscription) (*protocol.Response, error) {
	c.Mll.MRequest <- core.RingRequest{Opt: SYNC_SUB_OPT, Source: core.RingSync{Sub: core.Subscription{Type: in.Opt, NodeId: in.NodeId, Tag: in.Tag, Topic: in.Name, Endpoint: c.rpcEndpoint}}}
	return &protocol.Response{Successful: true, Message: "topic created"}, nil
}
func (c *DataServiceProvider) Unsubscribe(ctx context.Context, in *protocol.Subscription) (*protocol.Response, error) {
	c.Mll.MRequest <- core.RingRequest{Opt: SYNC_SUB_OPT, Source: core.RingSync{Sub: core.Subscription{Type: in.Opt, NodeId: in.NodeId, Tag: in.Tag, Topic: in.Name, Endpoint: c.rpcEndpoint, Deleting: true}}}
	return &protocol.Response{Successful: true, Message: "topic removed"}, nil
}

func (c *DataServiceProvider) Request(ctx context.Context, request *protocol.Request) (*protocol.Response, error) {
	switch request.Opt {

	case core.GET_DATA_REQUEST:
		return c.runGet(request)

	case core.CREATE_DATA_REQUEST:
		return c.runCreate(request)

	case core.UPDATE_DATA_REQUEST:
		return c.runUpdate(request)

	case core.DELETE_DATA_REQUEST:
		return c.runDelete(request)

	case core.RESET_DATA_REQUEST:
		return c.runReset(request)

	default:
	}
	return &protocol.Response{Successful: false, Message: "not suppotred"}, fmt.Errorf("opt not supported %d", request.Opt)
}

func (c *DataServiceProvider) List(in *protocol.Request, stream grpc.ServerStreamingServer[protocol.Response]) error {
	c.runQuery(in, stream)
	return nil
}

func (c *DataServiceProvider) Issue(ctx context.Context, task *protocol.Task) (*protocol.Response, error) {
	task.Meta.Id = c.tid()
	task.Meta.Prefix = c.tprefix(task.Meta.Id)
	task.Meta.Timeout = TASK_TIMEOUT_SECONDS
	task.Meta.Retries = TASK_RETRY_MAX
	for _, t := range task.Transactions {
		t.Meta.TaskId = task.Meta.Id
		t.Meta.Id = c.tid()
		t.Meta.State = protocol.TCC_RESERVING
		t.Meta.Timeout = TRANSACTION_TIMEOUT_SECONDS
		t.Meta.Retries = TASK_RETRY_MAX
	}
	tb := persistence.TaskBuilder{Target: task}
	req, err := tb.Request()
	if err != nil {
		return &protocol.Response{Successful: false, Message: err.Error()}, err
	}
	resp, err := c.runCreate(req)
	if err != nil {
		return resp, err
	}
	return c.runSetup(task)
}

func (c *DataServiceProvider) Confirm(ctx context.Context, meta *protocol.Meta) (*protocol.Response, error) {
	//call Confirmed
	c.runConfirmed(meta)
	return &protocol.Response{Successful: true}, nil
}

func (c *DataServiceProvider) Cancel(ctx context.Context, meta *protocol.Meta) (*protocol.Response, error) {
	//call Canceled
	c.runCanceled(meta)
	return &protocol.Response{Successful: true}, nil
}

func (c *DataServiceProvider) Finish(ctx context.Context, meta *protocol.Meta) (*protocol.Response, error) {
	//call Canceled
	c.runFinished(meta)
	return &protocol.Response{Successful: true}, nil
}

func (c *DataServiceProvider) tid() uint64 {
	for {
		id, err := c.seq.Id()
		if err == nil {
			return uint64(id)
		}
		time.Sleep(1 * time.Millisecond)
	}
}
func (c *DataServiceProvider) tprefix(id uint64) uint32 {
	buff := core.NewBuffer(8)
	buff.WriteUInt64(id)
	buff.Flip()
	bt, _ := buff.Read(0)
	return c.Mll.RingToken(bt)
}
