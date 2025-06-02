package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type PresenceRegister struct {
	*PresenceService
}

func (s *PresenceRegister) AccessControl() int32 {
	return bootstrap.PUBLIC_ACCESS_CONTROL
}

func (s *PresenceRegister) Register(login *event.Login) {
	id, _ := s.Seq.Id()
	login.SystemId = id
	login.AccessControl = bootstrap.PROTECTED_ACCESS_CONTROL
	hash, _ := s.Auth.HashPassword(login.Hash)
	login.Hash = hash
	err := s.SaveLogin(login)
	if err != nil {
		login.Cc <- event.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), DB_OP_ERR_CODE)}
		return
	}
	tk, err := s.Auth.CreateToken(login.SystemId, login.SystemId, login.AccessControl)
	if err != nil {
		login.Cc <- event.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), INVALID_TOKEN_CODE)}
		return
	}
	session := core.OnSession{Successful: true, SystemId: login.SystemId, Stub: login.SystemId, Token: tk, Home: s.Cls.Local().HttpEndpoint}
	login.Cc <- event.Chunk{Remaining: false, Data: util.ToJson(session)}
	s.Publish(login)
}

func (s *PresenceRegister) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	listener := make(chan event.Chunk)
	defer func() {
		close(listener)
		r.Body.Close()
	}()
	w.WriteHeader(http.StatusOK)
	var login event.Login
	json.NewDecoder(r.Body).Decode(&login)
	login.EventObj.Cc = listener
	go s.Register(&login)
	for c := range listener {
		w.Write(c.Data)
		if !c.Remaining {
			break
		}
	}

}
