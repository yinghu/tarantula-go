package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type AdminChangePwd struct {
	*AdminService
}

func (s *AdminChangePwd) AccessControl() int32 {
	return bootstrap.PROTECTED_ACCESS_CONTROL
}
func (s *AdminChangePwd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	hash, err := s.Auth.HashPassword(pwd)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	login.Hash = hash
	err = s.UpdatePassword(&login)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	session := core.OnSession{Successful: true, Message: "password changed"}
	w.Write(util.ToJson(session))
}
