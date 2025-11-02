package event

import (
	"encoding/json"
	"fmt"
	
	"testing"
	"time"

	"gameclustering.com/internal/core"
)

type SampleCreator struct {
}

func (s *SampleCreator) Create(cid int, topic string) (Event, error) {
	e := CreateEvent(cid)
	e.OnTopic(topic)
	e.OnListener(s)
	return e, nil

}
func (s *SampleCreator) VerifyTicket(ticket string) (core.OnSession, error) {
	sess := core.OnSession{Successful: true,SystemId: 100}
	return sess, nil
}
func (s *SampleCreator) OnError(e Event, err error) {
	fmt.Printf("On event error %v %s\n", e, err.Error())
}

func (s *SampleCreator) OnEvent(e Event) {
	fmt.Printf("On event %v\n", e)

}
func (s *SampleCreator) Send(e Event) error {
	return nil
}
func (s *SampleCreator) List(q Query)    {}
func (s *SampleCreator) Recover(q Query) {}
func (s *SampleCreator) Load(e Query)    {}

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
	epp := TcpEndpoint{Endpoint: "tcp://localhost:5050", OutboundEnabled: true}
	epp.Service = &SampleCreator{}
	go epp.Open()
	time.Sleep(10 * time.Second)
	cpp := TcpPublisher{Remote: "tcp://localhost:5050"}
	cpp.Connect()
	je := JoinEvent{Ticket: "xticket"}
	je.OnTopic("tpp")
	cpp.Join(&je)
	time.Sleep(1 * time.Second)
	cpp.Publish(&je, "yppp")
	time.Sleep(5 * time.Second)
	cpp.client.Close()
	time.Sleep(5 * time.Second)
	epp.listener.Close()
}
