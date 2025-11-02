package event

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"gameclustering.com/internal/core"
)

func createEvent() Event {
	sub := SubscriptionEvent{}
	return &sub
}
func TestEventJson(t *testing.T) {
	core.CreateTestLog()
	sub := SubscriptionEvent{App: "presence", Name: "ban"}
	data, err := json.Marshal(sub)
	if err != nil {
		t.Errorf("json error %s", err.Error())
	}
	e := createEvent()
	json.Unmarshal(data, e)

	pb, isSub := e.(*SubscriptionEvent)
	if !isSub {
		t.Errorf("json error %s", err.Error())
	}
	if pb.App != sub.App {
		t.Errorf("app not there %s", err.Error())
	}

	fm := fmt.Sprintf(`"%s":%d`, "count", 10)
	fmt.Printf("format :%s\n", fm)
	bm := fmt.Appendf(make([]byte, 0), `"%s":%d`, "count", 10)
	fmt.Printf("format :%s\n", string(bm))
	epp := TcpEndpoint{Endpoint: "tcp://localhost:5050",OutboundEnabled: true}
	go epp.Open()
	time.Sleep(10*time.Second)
	cpp := TcpPublisher{Remote: "tcp://localhost:5050"}
	cpp.Connect()
	je := JoinEvent{Ticket: "xticket"}
	cpp.Join(&je)
	time.Sleep(5*time.Second)
	cpp.client.Close()
	epp.listener.Close()
}
