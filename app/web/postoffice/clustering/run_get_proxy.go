package clustering

import (
	context "context"
	"io"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
)

func (c *DataServiceProvider) runGet(set *protocol.Request, responser grpc.ServerStreamingServer[protocol.Response]) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	ch := make(chan *protocol.Response, 3)
	defer close(ch)
	kh := c.Mll.RingToken(set.Data.Key)
	c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: kh, Replicas: REPLICA_MAX, Async: rq}
	nodes := <-rq
	ringNode := nodes[0]
	c.WTask <- Task{target: ringNode.RpcEndpoint, execute: func(tcp *grpc.ClientConn, opt int) error {
		dsp := protocol.NewDataServiceClient(tcp)
		stream, err := dsp.Get(context.Background(), set)
		if err != nil {
			ch <- &protocol.Response{Successful: false, Message: err.Error(), Code: 500000}
			return nil
		}
		for {
			data, err := stream.Recv()
			if err == io.EOF {
				ch <- &protocol.Response{Successful: false, Message: err.Error(), Code: 500000}
				break
			}
			if err != nil {
				core.AppLog.Debug().Msgf("run get streaming error %s", err.Error())
				ch <- &protocol.Response{Successful: false, Message: err.Error(), Code: 500000}
				break
			}
			ch <- data
			if !data.Successful {
				break
			}
		}
		return nil
	}}
	for resp := range ch {
		responser.Send(resp)
		if !resp.Successful {
			break
		}
	}

}

func (c *DataServiceProvider) runPull(target string, set *protocol.Request, ch chan *protocol.Response) {
	core.AppLog.Debug().Msgf("run remote pull %s >= %d < %d", target, set.Prefix, set.Opt)
	c.WTask <- Task{target: target, execute: func(tcp *grpc.ClientConn, opt int) error {
		dsp := protocol.NewDataServiceClient(tcp)
		stream, err := dsp.Pull(context.Background(), set)
		if err != nil {
			ch <- &protocol.Response{Successful: false, Message: err.Error()}
			return nil
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
		return nil
	}}
}
