package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/event"
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
	hash, _ := s.Auth.HashPassword(login.Hash)
	login.Hash = hash
	err := s.SaveLogin(login)
	if err != nil {
		login.Cc <- event.Chunk{Remaining: false, Data: errorMessage(err.Error(), DB_OP_ERR_CODE)}
		return
	}
	login.Cc <- event.Chunk{Remaining: false, Data: successMessage("registered")}
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
