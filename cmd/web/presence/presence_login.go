package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type PresenceLogin struct {
	*PresenceService
}

func (s *PresenceLogin) AccessControl() int32 {
	return bootstrap.PUBLIC_ACCESS_CONTROL
}

func (s *PresenceLogin) Login(login *event.Login) {
	pwd := login.Hash
	err := s.LoadLogin(login)
	if err != nil {
		login.Cc <- event.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), DB_OP_ERR_CODE)}
		return
	}
	err = util.ValidatePassword(pwd, login.Hash)
	if err != nil {
		login.Cc <- event.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), WRONG_PASS_CODE)}
		return
	}
	tk, err := s.Auth.CreateToken(login.SystemId, login.SystemId, login.AccessControl)
	if err != nil {
		login.Cc <- event.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), INVALID_TOKEN_CODE)}
		return
	}
	session := core.OnSession{Successful: true, SystemId: login.SystemId, Stub: login.SystemId, Token: tk, Home: s.Cls.Local().HttpEndpoint}
	login.Cc <- event.Chunk{Remaining: false, Data: util.ToJson(session)}
}

func (s *PresenceLogin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	listener := make(chan event.Chunk)
	defer func() {
		close(listener)
		r.Body.Close()
	}()
	w.WriteHeader(http.StatusOK)
	var login event.Login
	json.NewDecoder(r.Body).Decode(&login)
	login.EventObj.Cc = listener
	go s.Login(&login)
	for c := range listener {
		w.Write(c.Data)
		if !c.Remaining {
			break
		}
	}

}
