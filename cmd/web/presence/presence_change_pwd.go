package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
)

type PresenceChangePwd struct {
	*PresenceService
}

func (s *PresenceChangePwd) chnagePwd(login bootstrap.Login) {
	pwd := login.Hash
	err := s.LoadLogin(&login)
	if err != nil {
		login.Cc <- core.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), bootstrap.DB_OP_ERR_CODE)}
		return
	}
	hash, err := s.Authenticator().HashPassword(pwd)
	if err != nil {
		login.Cc <- core.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), bootstrap.DB_OP_ERR_CODE)}
		return
	}
	login.Hash = hash
	err = s.UpdatePassword(&login)
	if err != nil {
		login.Cc <- core.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), bootstrap.DB_OP_ERR_CODE)}
		return
	}
	login.Cc <- core.Chunk{Remaining: false, Data: bootstrap.SuccessMessage("password changed")}
}

func (s *PresenceChangePwd) AccessControl() int32 {
	return bootstrap.PROTECTED_ACCESS_CONTROL
}

func (s *PresenceChangePwd) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	listener := make(chan core.Chunk)
	defer func() {
		close(listener)
		r.Body.Close()
	}()
	w.WriteHeader(http.StatusOK)
	var login bootstrap.Login
	json.NewDecoder(r.Body).Decode(&login)
	login.Cc = listener
	go s.chnagePwd(login)
	for c := range listener {
		w.Write(c.Data)
		if !c.Remaining {
			break
		}
	}

}
