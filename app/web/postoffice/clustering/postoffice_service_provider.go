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
		c.runGet(request, rc)
		for resp := range rc {
			stream.Send(resp)
			if !resp.Successful {
				break
			}
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
	default:
		stream.Send(&protocol.Response{Successful: false, Message: fmt.Sprintf("opt not suuported %d", request.Opt)})
	}
	return nil
}

func (c *DataServiceProvider) runCreate(set *protocol.Request, ch chan *protocol.Response) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: c.Mll.RingToken(set.Data.Key), Replicas: REPLICA_MAX, Async: rq}
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
	return dsp.Create(context.Background(), request)
}

func (c *DataServiceProvider) runDelete(set *protocol.Request, ch chan *protocol.Response) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: c.Mll.RingToken(set.Data.Key), Replicas: REPLICA_MAX, Async: rq}
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
	return dsp.Create(context.Background(), request)
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

func (m *DataServiceProvider) clientReset(target *core.Node, request *protocol.Request) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	return dsp.Create(context.Background(), request)
}

func (c *DataServiceProvider) runGet(set *protocol.Request, ch chan *protocol.Response) {
	core.AppLog.Debug().Msg("running get opt")
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: c.Mll.RingToken(set.Data.Key), Replicas: REPLICA_MAX, Async: rq}
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
		for {
			data, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				core.AppLog.Debug().Msgf("streaming error %s", err.Error())
				break
			}
			core.AppLog.Debug().Msgf("Rev : %v", data)
			ch <- data
		}
		retry.Suc = true
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err.Error(), retry.Reties)
	ch <- &protocol.Response{Successful: false, Message: retry.Err.Error()}
}
