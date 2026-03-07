package clustering

import (
	context "context"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
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

func (c *DataServiceProvider) Request(ctx context.Context, in *protocol.Request) (*protocol.Data, error) {
	//getdata := GetData{in}
	return nil, nil //c.get(getdata)
}
