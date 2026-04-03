package clustering

import (
	context "context"
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
)

func (c *DataServiceProvider) runCreate(set *protocol.Request) (*protocol.Response, error) {
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
	core.AppLog.Debug().Msgf("data header %v", set.Data.Header)
	ch := make(chan *protocol.Response, 3)
	defer close(ch)
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: rt, Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		c.clientCreate(&ringNode, set, ch)
		resp := <-ch
		if !resp.Successful {
			retry.Err = resp.Message
			retry.Reties--
			continue
		}
		retry.Suc = true
		slaves := nodes[1:]
		for _, slave := range slaves {
			c.clientCreate(&slave, set, ch)
			resp = <-ch
			if !resp.Successful {
				core.AppLog.Debug().Msgf("error on slave %s", resp.Message)
			}
		}
		break
	}
	core.AppLog.Printf("retry %s, %d %v", retry.Err, retry.Reties, retry.Suc)
	return &protocol.Response{Successful: retry.Suc, Message: retry.Err}, nil
}

func (m *DataServiceProvider) clientCreate(target *core.Node, request *protocol.Request, ch chan *protocol.Response) {
	task := core.Task{Target: target.RpcEndpoint, Execute: func(tcp *grpc.ClientConn, opt int) error {
		if opt == core.NO_TCP_CONNECT {
			ch <- &protocol.Response{Successful: false, Message: "no tcp connect"}
			return fmt.Errorf("no tcp connect")
		}
		dsp := protocol.NewDataServiceClient(tcp)
		resp, err := dsp.Create(context.Background(), request)
		if err != nil {
			ch <- &protocol.Response{Successful: false, Message: err.Error()}
			return err
		}
		ch <- resp
		return nil
	}}
	m.WTask <- task
}
