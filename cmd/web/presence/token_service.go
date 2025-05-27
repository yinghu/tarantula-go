package main

import (
	"errors"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

func (s *PresenceService) VerifyToken(token string, listener chan event.Chunk) {
	err := s.Tkn.Verify(token, func(h *core.JwtHeader, p *core.JwtPayload) error {
		t := time.UnixMilli(p.Exp).UTC()
		if t.Before(time.Now().UTC()) {
			return errors.New("token expired")
		}
		return nil
	})
	if err != nil {
		listener <- event.Chunk{Remaining: false, Data: errorMessage(err.Error(), INVALID_TOKEN_CODE)}
		return
	}
	listener <- event.Chunk{Remaining: false, Data: successMessage("passed")}
}
