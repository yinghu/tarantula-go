package main

import (
	"encoding/json"
	"net/http"
	"time"

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

func (s *PresenceLogin) Login(login bootstrap.Login) {
	pwd := login.Hash
	err := s.LoadLogin(&login)
	if err != nil {
		login.Cc <- event.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), bootstrap.DB_OP_ERR_CODE)}
		return
	}
	err = util.ValidatePassword(pwd, login.Hash)
	if err != nil {
		login.Cc <- event.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), bootstrap.WRONG_PASS_CODE)}
		return
	}
	tk, err := s.Authenticator().CreateToken(login.SystemId, login.Id, login.AccessControl)
	if err != nil {
		login.Cc <- event.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), bootstrap.INVALID_TOKEN_CODE)}
		return
	}
	session := core.OnSession{Successful: true, SystemId: login.SystemId, Stub: login.Id, Token: tk, Home: s.Cluster().Local().HttpEndpoint}
	login.Cc <- event.Chunk{Remaining: false, Data: util.ToJson(session)}
	go func() {
		id, err := s.Sequence().Id()
		if err != nil {
			return
		}
		me := event.LoginEvent{SystemId: login.SystemId, Name: login.Name}
		me.OnOId(id)
		me.LoginTime = time.Now()
		me.Topic("login")
		s.Send(&me)
	}()
}

func (s *PresenceLogin) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	listener := make(chan event.Chunk)
	defer func() {
		close(listener)
		r.Body.Close()
	}()
	w.WriteHeader(http.StatusOK)
	var login bootstrap.Login
	json.NewDecoder(r.Body).Decode(&login)
	login.Cc = listener
	go s.Login(login)
	for c := range listener {
		w.Write(c.Data)
		if !c.Remaining {
			break
		}
	}

}
