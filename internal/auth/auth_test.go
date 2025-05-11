package auth

import (
	//"fmt"
	"testing"
)

func TestAuth(t *testing.T) {
	service := Service{NodeId: 1, DatabaseURL: "postgres://postgres:password@192.168.1.7:5432/tarantula_user"}
	err := service.Start()
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
