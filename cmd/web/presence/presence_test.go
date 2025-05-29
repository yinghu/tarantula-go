package main

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
	"gameclustering.com/internal/core"
)

func TestAuth(t *testing.T) {
	service := PresenceService{}
	f := conf.Env{}
	f.Load("/etc/tarantula/presence-conf.json")
	err := service.Start(f, nil)
	if err != nil {
		t.Errorf("Service error %s", err.Error())
	}
	defer service.Shutdown()
	login := event.Login{Name: "foo1003", Hash: "ppp", SystemId: 10, ReferenceId: 1}
	login.Cc = make(chan event.Chunk)
	go service.Register(&login)
	for c := range login.Cc {
		if !c.Remaining {
			break
		}
	}
	ses := core.OnSession{Message: WRONG_PASS_MSG, ErrorCode: WRONG_PASS_CODE}
	ser := util.ToJson(ses)
	fmt.Printf("JSON :%s\n", string(ser))
}
