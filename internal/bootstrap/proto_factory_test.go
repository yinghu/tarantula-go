package bootstrap

import (
	"fmt"
	"testing"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMessageEventFactory(t *testing.T) {

	me := protocol.MessageEvent{Title: "tile", Message: "msg", DateTime: timestamppb.New(time.Now()), Source: "admin"}

	ptf := MessageEventFactory{}
	tp, err := ptf.FromMessageEvent(&me)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	tp.NodeId = "nodeId"
	tp.Tag = "presence"
	tp.Name = "message"
	tp.Event.Id = 100
	req, err := ptf.Request(tp)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	fmt.Printf("req %v\n", &req.Data.Value)
	tpc, err := ptf.Topic(req.Data.Value)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	fmt.Printf("tp %v\n", tp)
	fmt.Printf("tpx %v\n", tpc)

	mq := MessageEventQuery{}
	mq.ClassId = core.EVENT_FACTORY_ID
	mq.FactoryId = MESSAGE_EVENT_CID
	mq.Topic = "message"

	qt, _ := ptf.Export(&mq)
	q, _ := ptf.Import(qt)
	fmt.Printf("q %d %d %s\n", q.QFactoryId(), q.QClassId(), mq.QTopic())
}
