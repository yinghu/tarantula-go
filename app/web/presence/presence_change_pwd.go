package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

type PresenceChangePwd struct {
	*PresenceService
}

func (s *PresenceChangePwd) chnagePwd(login *protocol.LoginObject) {
	pwd := login.Password
	err := s.LoadLogin(login)
	if err != nil {
		//login.Cc <- core.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), bootstrap.DB_OP_ERR_CODE)}
		return
	}
	hash, err := s.Authenticator().HashPassword(pwd)
	if err != nil {
		//login.Cc <- core.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), bootstrap.DB_OP_ERR_CODE)}
		return
	}
	login.Password = hash
	err = s.UpdatePassword(login)
	if err != nil {
		//login.Cc <- core.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), bootstrap.DB_OP_ERR_CODE)}
		return
	}
	//login.Cc <- core.Chunk{Remaining: false, Data: bootstrap.SuccessMessage("password changed")}
}

func (s *PresenceChangePwd) AccessControl() int32 {
	return core.PROTECTED_ACCESS_CONTROL
}

func (s *PresenceChangePwd) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	listener := make(chan core.Chunk)
	defer func() {
		close(listener)
		r.Body.Close()
	}()
	w.WriteHeader(http.StatusOK)
	var login protocol.LoginObject
	json.NewDecoder(r.Body).Decode(&login)
	//login.Cc = listener
	go s.chnagePwd(&login)
	for c := range listener {
		cv, _ := c.Data.([]byte)
		w.Write(cv)
		if !c.Remaining {
			break
		}
	}

}
