package auth

import (
	"testing"

	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/event"
)

func TestAuth(t *testing.T) {
	service := Service{}
	f := conf.Env{}
	f.Load("/etc/tarantula/presence-conf.json")
	err := service.Start(f)
	if err != nil {
		t.Errorf("Service error %s", err.Error())
	}
	defer service.Shutdown()
	login := Login{Name: "foo1003", Hash: "ppp", SystemId: 10, ReferenceId: 1}
	login.Listener = make(chan event.Chunk)
	login.Topic = false
	if login.OnTopic() {
		t.Errorf("login topic error %v", login.OnTopic())
	}
	go service.Register(&login)
	for c := range login.Listener {
		if !c.Remaining {
			break
		}
	}
}
