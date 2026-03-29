package bootstrap

import (
	"fmt"
	"testing"
	"time"

	"gameclustering.com/internal/event"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"
)

func TestProtoFactory(t *testing.T) {
	me := event.MessageEvent{Title: "tile", Message: "msg", DateTime: time.Now(), Source: "admin"}
	me.OnNodeId("nid")
	me.OnTag("presence")
	me.OnTopic("message")
	me.OnOId(100)
	tp, err := FromMessageEvent(me)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	req, err := ToRequest(tp)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	fmt.Printf("req %v\n", req)
	tpx := protocol.Topic{}
	err = proto.Unmarshal(req.Data.Value, &tpx)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	fmt.Printf("tp %v\n", tp)
	fmt.Printf("tpx %v\n", &tpx)

}
