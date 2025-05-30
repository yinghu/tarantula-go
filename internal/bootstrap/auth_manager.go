package bootstrap

import (
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AuthManager struct {
	Tkn core.Jwt
}

func (s *AuthManager) HashPassword(password string) (string, error) {
	return util.HashPassword(password)
}
func (s *AuthManager) ValidatePassword(password string, hash string) error {
	return util.ValidatePassword(password, hash)
}
func (s *AuthManager) CreateToken(systemId int64, stub int64) (string, error) {
	return s.Tkn.Token(func(h *core.JwtHeader, p *core.JwtPayload) error {
		h.Kid = "kid"
		p.Aud = "player"
		exp := time.Now().Add(time.Hour * 24).UTC()
		p.Exp = exp.UnixMilli()
		return nil
	})
}
func (s *AuthManager) ValidateToken(token string) error {
	return s.Tkn.Verify(token,func(jh *core.JwtHeader, jp *core.JwtPayload) error {
		
		return nil
	})
}
