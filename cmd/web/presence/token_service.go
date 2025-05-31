package main

import (
	"gameclustering.com/internal/event"
)

func (s *PresenceService) VerifyToken(token string, listener chan event.Chunk) {
	_ , err := s.Auth.ValidateToken(token)
	if err != nil {
		listener <- event.Chunk{Remaining: false, Data: errorMessage(err.Error(), INVALID_TOKEN_CODE)}
		return
	}
	listener <- event.Chunk{Remaining: false, Data: successMessage("passed")}
}
