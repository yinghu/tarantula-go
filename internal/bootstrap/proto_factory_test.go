package bootstrap

import (
	"fmt"
	"testing"
	"time"

	"gameclustering.com/internal/event"
)

func TestProtoFactory(t *testing.T) {
	me := event.MessageEvent{Title: "tile", Message: "msg", DateTime: time.Now(), Source: "admin"}
	me.OnNodeId("nid")
	me.OnTag("presence")
	me.OnTopic("message")
	me.OnOId(100)
	ptf := MessageEventFactory{}
	ptf.Cluster = &ClusterManager{}
	tp, err := ptf.FromMessageEvent(me)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	req, err := ptf.ToRequest(tp)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	fmt.Printf("req %v\n", req)
	tpx, err := ptf.FromData(req.Data.Value)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	fmt.Printf("tp %v\n", tp)
	fmt.Printf("tpx %v\n", tpx)

}
