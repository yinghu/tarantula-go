package main

import (
	"testing"
	"time"

	"gameclustering.com/internal/event"
)

type SampleCreator struct {
}

func (s *SampleCreator) Create(cid int, topic string) (event.Event, error) {
	me := MahjongEvent{}
	me.Callback = &MahjongEventListener{}
	return &me, nil
}

func TestClient(t *testing.T) {
	sb := event.SocketPublisher{Remote: "tcp://192.168.1.11:5050"}
	err := sb.Connect()
	if err != nil {
		t.Errorf("conn error %s", err.Error())
	}

	e := event.JoinEvent{Token: "test token"}
	e.OnListener(&MahjongEventListener{})
	err = sb.Join(&e)
	if err != nil {
		t.Errorf("send error %s", err.Error())
		sb.Close()
		return
	}
	go sb.Subscribe(&SampleCreator{}, &MahjongEventListener{})
	for range 3 {
		me := MahjongEvent{Cmd: "drop"}
		me.OnListener(&MahjongEventListener{})
		sb.Publish(&me, "validated")
	}
	time.Sleep(5 * time.Second)
	sb.Close()
}
