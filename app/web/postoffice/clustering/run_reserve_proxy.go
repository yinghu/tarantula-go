package clustering

import (
	context "context"
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

func (c *DataServiceProvider) runReserve(t *protocol.Transaction) (*protocol.Response, error) {
	rq := make(chan []core.Subscription, 3)
	defer close(rq)
	c.DRequest <- TopicRequest{Opt: TOPIC_REGISTER, Subs: rq, NodeId: t.Meta.NodeId, Tag: t.Meta.Tag, Name: fmt.Sprintf("%s%s", TRANS_SUB_PREFIX, t.Meta.Name)}
	subs := <-rq
	for _, sub := range subs {
		conn, err := sub.CPool.Conn()
		if err != nil {
			core.AppLog.Warn().Msgf("no connection available on sub %v", sub)
			continue
		}
		dsp := protocol.NewTransactionServiceClient(conn.Conn)
		return dsp.Reserve(context.Background(), t)
	}
	//rq := make(chan []core.Node, 3)
	//defer close(rq)

	//c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: t.Meta.Prefix, Replicas: REPLICA_MAX, Async: rq}
	//nodes := <-rq
	//ringNode := nodes[0]
	//conn, err := ringNode.CPool.Conn()

	//if err != nil {
	return &protocol.Response{Successful: false}, fmt.Errorf("no subscription available")
	//}
	//dsp := protocol.NewTransactionServiceClient(conn.Conn)
	//return dsp.Reserve(context.Background(), t)
}
