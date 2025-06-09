package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type AdminLogin struct {
	*AdminService
}

func (s *AdminLogin) AccessControl() int32 {
	return bootstrap.PUBLIC_ACCESS_CONTROL
}
func (s *AdminLogin) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var login event.Login
	json.NewDecoder(r.Body).Decode(&login)
	pwd := login.Hash
	err := s.LoadLogin(&login)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	err = s.Authenticator().ValidatePassword(pwd, login.Hash)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	tk, err := s.Authenticator().CreateToken(login.SystemId, login.Id, login.AccessControl)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	session := core.OnSession{Successful: true, SystemId: login.SystemId, Stub: login.Id, Token: tk, Home: s.Cluster().Local().HttpEndpoint}
	w.Write(util.ToJson(session))
}
