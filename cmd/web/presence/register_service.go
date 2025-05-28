package main

import (
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

func (s *PresenceService) Register(login *event.Login) {
	id, _ := s.Seq.Id()
	login.SystemId = id
	hash, _ := util.HashPassword(login.Hash)
	login.Hash = hash
	err := s.SaveLogin(login)
	if err != nil {
		login.Cc <- event.Chunk{Remaining: false, Data: errorMessage(err.Error(), DB_OP_ERR_CODE)}
		return
	}
	login.Cc <- event.Chunk{Remaining: false, Data: successMessage("registered")}
	s.Publish(login)
}




