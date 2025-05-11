package auth

import (
	//"fmt"
	"testing"

	"gameclustering.com/internal/conf"
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
	login := Login{Name: "foo100", Hash: "ppp", SystemId: 2, ReferenceId: 1}
	er := service.Register(&login)
	if er != nil {
		t.Errorf("Register error %s", er.Error())
	}
}
