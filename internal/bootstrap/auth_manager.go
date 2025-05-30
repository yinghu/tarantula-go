package bootstrap

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AuthManager struct {
	Kid string
	Tkn core.Jwt
	Cip *util.Cipher
}

func (s *AuthManager) HashPassword(password string) (string, error) {
	return util.HashPassword(password)
}
func (s *AuthManager) ValidatePassword(password string, hash string) error {
	return util.ValidatePassword(password, hash)
}
func (s *AuthManager) CreateToken(systemId int64, stub int64) (string, error) {
	return s.Tkn.Token(func(h *core.JwtHeader, p *core.JwtPayload) error {
		h.Kid = s.Kid
		exp := time.Now().Add(time.Hour * 24).UTC()
		p.Exp = exp.UnixMilli()
		aud := fmt.Sprintf("%d.%d.%d", systemId, stub, p.Exp)
		p.Aud = s.Cip.Encrypt(aud)
		fmt.Printf("%s\n", p.Aud)
		return nil
	})
}
func (s *AuthManager) ValidateToken(token string) error {
	return s.Tkn.Verify(token, func(jh *core.JwtHeader, jp *core.JwtPayload) error {
		aud, err := s.Cip.Decrypt(jp.Aud)
		if err != nil {
			return err
		}
		parts := strings.Split(aud, ".")
		sysId, err := strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return err
		}
		stub, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			return err
		}
		exp, err := strconv.ParseUint(parts[2], 10, 64)
		if err != nil {
			return err
		}
		fmt.Printf("%s, %s, %d, %d, %d\n", aud, jh.Kid, sysId, stub, exp)
		return nil
	})
}
