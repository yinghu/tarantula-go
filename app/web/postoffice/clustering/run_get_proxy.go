package clustering

import (
	context "context"
	"io"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (c *DataServiceProvider) runGet(set *protocol.Request, responser grpc.ServerStreamingServer[protocol.Response]) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	kh := c.Mll.RingToken(set.Data.Key)
	c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: kh, Replicas: REPLICA_MAX, Async: rq}
	nodes := <-rq
	ringNode := nodes[0]
	conn, err := ringNode.CPool.Conn()
	if err != nil {
		responser.Send(&protocol.Response{Successful: false, Message: err.Error()})
		return
	}
	dsp := protocol.NewDataServiceClient(conn.Conn)
	stream, err := dsp.Get(context.Background(), set)
	if err != nil {
		responser.Send(&protocol.Response{Successful: false, Message: err.Error(), Code: 500000})
		return
	}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			responser.Send(&protocol.Response{Successful: false, Message: err.Error(), Code: 500000})
			break
		}
		if err != nil {
			core.AppLog.Debug().Msgf("run get streaming error %s", err.Error())
			responser.Send(&protocol.Response{Successful: false, Message: err.Error(), Code: 500000})
			break
		}
		responser.Send(data)
		if !data.Successful {
			break
		}
	}

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
