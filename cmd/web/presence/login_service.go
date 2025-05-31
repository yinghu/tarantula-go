package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

func (s *PresenceService) Login(login *event.Login) {
	pwd := login.Hash
	err := s.LoadLogin(login)
	if err != nil {
		login.Cc <- event.Chunk{Remaining: false, Data: errorMessage(err.Error(), DB_OP_ERR_CODE)}
		return
	}
	err = util.ValidatePassword(pwd, login.Hash)
	if err != nil {
		login.Cc <- event.Chunk{Remaining: false, Data: errorMessage(err.Error(), WRONG_PASS_CODE)}
		return
	}
	tk, err := s.Auth.CreateToken(login.SystemId, login.SystemId,0)
	if err != nil {
		login.Cc <- event.Chunk{Remaining: false, Data: errorMessage(err.Error(), INVALID_TOKEN_CODE)}
		return
	}
	session := core.OnSession{Successful: true, SystemId: login.SystemId, Stub: login.SystemId, Token: tk, Home: s.Cluster.Local().HttpEndpoint}
	login.Cc <- event.Chunk{Remaining: false, Data: util.ToJson(session)}
}
