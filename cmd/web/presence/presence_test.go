package main

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
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
	reg := PresenceRegister{&service}
	go reg.Register(&login)
	for c := range login.Cc {
		if !c.Remaining {
			break
		}
	}
	ses := core.OnSession{Message: bootstrap.WRONG_PASS_MSG, ErrorCode: bootstrap.WRONG_PASS_CODE}
	ser := util.ToJson(ses)
	fmt.Printf("JSON :%s\n", string(ser))
}
