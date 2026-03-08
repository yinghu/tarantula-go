package clustering

import (
	context "context"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (c *DataServiceProvider) HashRing(request *protocol.Request, stream grpc.ServerStreamingServer[protocol.HashNode]) error {
	rq := make(chan []core.Node, 1)
	defer close(rq)
	c.Mll.rangeRing(core.RingRequest{Async: rq, Opt: core.ALL_RING_OPT})
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
	c.Mll.rangeRing(core.RingRequest{Async: rq, Opt: core.REPLICA_RING_OPT, Token: request.Prefix})
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
	core.AppLog.Debug().Msgf("requesting data %v", request)
	rc := make(chan *protocol.Response, 3)
	defer close(rc)
	c.runCreate(request, rc)
	resp := <-rc
	stream.Send(resp)
	return nil
}

func (c *DataServiceProvider) runCreate(set *protocol.Request, ch chan *protocol.Response) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Token: c.Mll.RingToken(set.Data.Key), Replicas: REPLICA_MAX, Async: rq}
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
