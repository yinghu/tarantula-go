package main

import (
	"encoding/json"
	"net/http"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type AdminLogin struct {
	*AdminService
}

func (s AdminLogin) Login(login *event.Login) error {

	return nil
}

func (s *AdminLogin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	err = util.ValidatePassword(pwd, login.Hash)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	tk, err := s.Tkn.Token(func(h *core.JwtHeader, p *core.JwtPayload) error {
		h.Kid = "kid"
		p.Aud = "player"
		exp := time.Now().Add(time.Hour * 24).UTC()
		p.Exp = exp.UnixMilli()
		return nil
	})
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	session := core.OnSession{Successful: true, SystemId: login.SystemId, Stub: login.SystemId, Token: tk, Home: s.Cluster().Local().HttpEndpoint}
	w.Write(util.ToJson(session))
}
