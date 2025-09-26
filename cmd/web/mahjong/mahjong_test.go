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
	for range 3{
		e := MahjongEvent{Cmd: "drop"}
		e.OnListener(&MahjongEventListener{})
		err = sb.Publish(&e, "test")
		if err != nil {
			t.Errorf("send error %s", err.Error())
		}
	}
	sb.Close()
}
