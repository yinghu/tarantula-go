package clustering

import (
	context "context"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

func (c *DataServiceProvider) runUpdate(set *protocol.Request) (*protocol.Response, error) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	ch := make(chan *protocol.Response, 3)
	defer close(ch)
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: c.Mll.RingToken(set.Data.Key), Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		c.clientUpdate(&ringNode, set, ch)
		resp := <-ch
		if !resp.Successful {
			retry.Err = resp.Message
			retry.Reties--
			continue
		}
		retry.Suc = true
		slaves := nodes[1:]
		for _, slave := range slaves {
			c.clientUpdate(&slave, set, ch)
		}
		break
	}
	core.AppLog.Printf("retry %s, %d ,%v", retry.Err, retry.Reties, retry.Suc)
	return &protocol.Response{Successful: retry.Suc, Message: retry.Err}, nil
}

func (m *DataServiceProvider) clientUpdate(target *core.Node, request *protocol.Request, ch chan *protocol.Response) {
	conn, err := target.CPool.Conn()
	if err != nil {
		ch <- &protocol.Response{Successful: false, Message: err.Error()}
		return
	}
	dsp := protocol.NewDataServiceClient(conn.Conn)
	resp, _ := dsp.Update(context.Background(), request)
	ch <- resp

}
