package event

import (
	"encoding/json"
	"testing"
)

func createEvent() Event {
	sub := SubscriptionEvent{}
	return &sub
}
func TestEventJson(t *testing.T) {
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
}
