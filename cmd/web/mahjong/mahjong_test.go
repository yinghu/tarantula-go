package main

import (
	"testing"

	"gameclustering.com/internal/event"
)

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
	go sb.Subscribe(&MahjongEventListener{})
	for range 3 {
		me := MahjongEvent{Cmd: "drop"}
		me.OnListener(&MahjongEventListener{})
		sb.Publish(&me, "validated")
	}
	sb.Close()
}
