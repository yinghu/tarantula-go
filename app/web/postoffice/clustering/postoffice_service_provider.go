package clustering

import (
	context "context"
	"fmt"
	"io"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func (c *DataServiceProvider) Receive(topic *protocol.Topic, stream grpc.ServerStreamingServer[protocol.Response]) error {
	core.AppLog.Debug().Msg("starting publish")
	aq := make(chan chan *protocol.Response, 2)
	c.DRequest <- TopicRequest{Opt: RECEIVER_START, Name: topic.Prefix, Rev: aq}
	ch := <-aq
	close(aq)
	core.AppLog.Debug().Msgf("starting publish on [%s]", topic.Prefix)
	for c.running {
		for resp := range ch {
			core.AppLog.Debug().Msgf("distributing message %s", resp)
			err := stream.Send(resp)
			if err != nil {
				core.AppLog.Debug().Msgf("send error %s", err.Error())
				break
			}
		}
	}
	core.AppLog.Debug().Msg("stop publish")
	return nil
}

func (c *DataServiceProvider) Subscribe(ctx context.Context, in *protocol.Topic) (*protocol.Response, error) {
	c.Mll.MRequest <- core.RingRequest{Opt: SYNC_NODE_OPT, Source: core.RingSync{Sub: core.Subscription{Prefix: in.Prefix, Topic: in.Name, Endpoint: c.rpcEndpoint}}}
	return &protocol.Response{Successful: true, Message: "topic created"}, nil
}
func (c *DataServiceProvider) Unsubscribe(ctx context.Context, in *protocol.Topic) (*protocol.Response, error) {
	c.Mll.MRequest <- core.RingRequest{Opt: SYNC_NODE_OPT, Source: core.RingSync{Sub: core.Subscription{Prefix: in.Prefix, Topic: in.Name, Endpoint: c.rpcEndpoint}}}
	return &protocol.Response{Successful: true, Message: "topic removed"}, nil
}

func (c *DataServiceProvider) Request(request *protocol.Request, stream grpc.ServerStreamingServer[protocol.Response]) error {
	switch request.Opt {
	case core.CREATE_DATA_REQUEST:
		rc := make(chan *protocol.Response, 3)
		defer close(rc)
		c.runCreate(request, rc)
		resp := <-rc
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
	case core.PUBLISH_REQUEST:
		//rc := make(chan *protocol.Response, 3)
		//defer close(rc)
		//c.runCreate(request, rc)
		//resp := <-rc
		//if !resp.Successful {
		//	stream.Send(resp)
		//} else {
		resp := protocol.Response{Successful: true, Message: "message published"}
		c.DMessager <- &resp
		//c.runPublish(request, rc)
		core.AppLog.Debug().Msgf("messaging out %s", resp.Message)
		//resp = <-rc
		stream.Send(&resp)
		//}
	default:
		stream.Send(&protocol.Response{Successful: false, Message: fmt.Sprintf("request opt not suuported %d", request.Opt)})
	}
	return nil
}

func (c *DataServiceProvider) runCreate(set *protocol.Request, ch chan *protocol.Response) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	var rt uint32
	if set.Prefix > 0 {
		rt = set.Prefix
		core.AppLog.Debug().Msgf("using prefix %d", set.Prefix)
	} else {
		rt = c.Mll.RingToken(set.Data.Key)
	}
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: rt, Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		resp, err := c.clientCreate(&ringNode, set)
		if err != nil {
			retry.Err = err
			retry.Reties--
			continue
		}
		ch <- resp
		retry.Suc = true
		if !resp.Successful {
			break
		}
		slaves := nodes[1:]
		for _, slave := range slaves {
			c.clientCreate(&slave, set)
		}
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err.Error(), retry.Reties)
	ch <- &protocol.Response{Successful: false, Message: retry.Err.Error()}
}

func (m *DataServiceProvider) clientCreate(target *core.Node, request *protocol.Request) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	return dsp.Create(context.Background(), request)
}

func (c *DataServiceProvider) runUpdate(set *protocol.Request, ch chan *protocol.Response) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: c.Mll.RingToken(set.Data.Key), Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		resp, err := c.clientUpdate(&ringNode, set)
		if err != nil {
			retry.Err = err
			retry.Reties--
			continue
		}
		ch <- resp
		retry.Suc = true
		if !resp.Successful {
			break
		}
		slaves := nodes[1:]
		for _, slave := range slaves {
			c.clientUpdate(&slave, set)
		}
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err.Error(), retry.Reties)
	ch <- &protocol.Response{Successful: false, Message: retry.Err.Error()}
}

func (m *DataServiceProvider) clientUpdate(target *core.Node, request *protocol.Request) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	return dsp.Update(context.Background(), request)
}

func (c *DataServiceProvider) runDelete(set *protocol.Request, ch chan *protocol.Response) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: c.Mll.RingToken(set.Data.Key), Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		resp, err := c.clientDelete(&ringNode, set)
		if err != nil {
			retry.Err = err
			retry.Reties--
			continue
		}
		ch <- resp
		retry.Suc = true
		if !resp.Successful {
			break
		}
		slaves := nodes[1:]
		for _, slave := range slaves {
			c.clientDelete(&slave, set)
		}
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err.Error(), retry.Reties)
	ch <- &protocol.Response{Successful: false, Message: retry.Err.Error()}
}

func (m *DataServiceProvider) clientDelete(target *core.Node, request *protocol.Request) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	return dsp.Delete(context.Background(), request)
}

func (c *DataServiceProvider) runReset(set *protocol.Request, ch chan *protocol.Response) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: c.Mll.RingToken(set.Data.Key), Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		resp, err := c.clientReset(&ringNode, set)
		if err != nil {
			retry.Err = err
			retry.Reties--
			continue
		}
		ch <- resp
		retry.Suc = true
		if !resp.Successful {
			break
		}
		slaves := nodes[1:]
		for _, slave := range slaves {
			c.clientReset(&slave, set)
		}
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err.Error(), retry.Reties)
	ch <- &protocol.Response{Successful: false, Message: retry.Err.Error()}
}

func (m *DataServiceProvider) clientReset(target *core.Node, request *protocol.Request) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	return dsp.Reset(context.Background(), request)
}

func (c *DataServiceProvider) runGet(set *protocol.Request, ch chan *protocol.Response) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	for retry.Reties > 0 {
		kh := c.Mll.RingToken(set.Data.Key)
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: kh, Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		tcp, err := grpc.NewClient(ringNode.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			retry.Reties++
			continue
		}
		defer tcp.Close()
		dsp := protocol.NewDataServiceClient(tcp)
		stream, err := dsp.Get(context.Background(), set)
		if err != nil {
			retry.Reties++
			continue
		}
		crt := protocol.Response{Successful: false}
		for {
			data, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				core.AppLog.Debug().Msgf("run get streaming error %s", err.Error())
				crt.Code = 400001
				crt.Message = err.Error()
				break
			}
			ch <- data
			if !data.Successful {
				break
			}
		}
		ch <- &crt
		retry.Suc = true
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err.Error(), retry.Reties)
	ch <- &protocol.Response{Successful: false, Message: retry.Err.Error()}
}

func (c *DataServiceProvider) runPull(target string, set *protocol.Request, ch chan *protocol.Response) {
	core.AppLog.Debug().Msgf("run remote pull %s >= %d < %d", target, set.Prefix, set.Opt)
	tcp, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		ch <- &protocol.Response{Successful: false, Message: err.Error()}
		return
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	stream, err := dsp.Pull(context.Background(), set)
	if err != nil {
		ch <- &protocol.Response{Successful: false, Message: err.Error()}
		return
	}
	crt := protocol.Response{Successful: false}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			core.AppLog.Debug().Msgf("run pull streaming error %s", err.Error())
			crt.Message = err.Error()
			break
		}
		ch <- data
		if !data.Successful {
			break
		}
	}
	ch <- &crt
}

func (c *DataServiceProvider) runPublish(set *protocol.Request, ch chan *protocol.Response) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: c.Mll.RingToken(set.Data.Key), Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		resp, err := c.clientPublish(&ringNode, set)
		if err != nil {
			retry.Err = err
			retry.Reties--
			continue
		}
		ch <- resp
		retry.Suc = true
		if !resp.Successful {
			break
		}
		slaves := nodes[1:]
		for _, slave := range slaves {
			c.clientPublish(&slave, set)
		}
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err.Error(), retry.Reties)
	ch <- &protocol.Response{Successful: false, Message: retry.Err.Error()}
}
func (m *DataServiceProvider) clientPublish(target *core.Node, request *protocol.Request) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	return dsp.Publish(context.Background(), request)
}
