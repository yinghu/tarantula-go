package clustering

import (
	context "context"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

func (c *DataServiceProvider) runSetup(t *protocol.Task) (*protocol.Response, error) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	kh := t.Meta.Prefix
	c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: kh, Replicas: REPLICA_MAX, Async: rq}
	nodes := <-rq
	ringNode := nodes[0]
	conn, err := ringNode.CPool.Conn()
	if err != nil {
		return &protocol.Response{Successful: false, Message: err.Error()}, err
	}
	dsp := protocol.NewTransactionServiceClient(conn.Conn)
	return dsp.Setup(context.Background(), t)
}
