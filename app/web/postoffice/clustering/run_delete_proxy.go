package clustering

import (
	context "context"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
)

func (c *DataServiceProvider) runDelete(set *protocol.Request) (*protocol.Response, error) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	ch := make(chan *protocol.Response, 3)
	defer close(ch)
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: c.Mll.RingToken(set.Data.Key), Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		c.clientDelete(&ringNode, set, ch)
		resp := <-ch
		if !resp.Successful {
			retry.Err = resp.Message
			retry.Reties--
			continue
		}
		retry.Suc = true
		slaves := nodes[1:]
		for _, slave := range slaves {
			c.clientDelete(&slave, set, ch)
		}
		break
	}
	core.AppLog.Printf("retry %s, %d , %v", retry.Err, retry.Reties, retry.Suc)
	return &protocol.Response{Successful: retry.Suc, Message: retry.Err}, nil
}

func (m *DataServiceProvider) clientDelete(target *core.Node, request *protocol.Request, ch chan *protocol.Response) {
	task := core.Task{Target: target.RpcEndpoint, Execute: func(tcp *grpc.ClientConn, opt int) error {
		dsp := protocol.NewDataServiceClient(tcp)
		resp, _ := dsp.Delete(context.Background(), request)
		ch <- resp
		return nil
	}}
	m.WTask <- task
}
