package bootstrap

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

const TOKEN_TIME_OUT_HOURS int = 24

type AuthManager struct {
	
	Kid      string
	Tkn      core.Jwt
	Cipher   *util.Aes
}

func (s *AuthManager) HashPassword(password string) (string, error) {
	return util.HashPassword(password)
}
func (s *AuthManager) ValidatePassword(password string, hash string) error {
	return util.ValidatePassword(password, hash)
}
func (s *AuthManager) CreateToken(systemId int64, stub int64, accessControl int32) (string, error) {
	return s.Tkn.Token(func(h *core.JwtHeader, p *core.JwtPayload) error {
		h.Kid = s.Kid
		exp := time.Now().Add(time.Hour * time.Duration(TOKEN_TIME_OUT_HOURS)).UTC()
		p.Exp = exp.UnixMilli()
		aud := fmt.Sprintf("%d.%d.%d.%d", systemId, stub, accessControl, p.Exp)
		p.Aud = s.Cipher.Encrypt(aud)
		return nil
	})
}
func (s *AuthManager) ValidateToken(token string) (core.OnSession, error) {
	session := core.OnSession{Successful: false}
	err := s.Tkn.Verify(token, func(jh *core.JwtHeader, jp *core.JwtPayload) error {
		aud, err := s.Cipher.Decrypt(jp.Aud)
		if err != nil {
			return err
		}
		parts := strings.Split(aud, ".")
		sysId, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return err
		}
		stub, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return err
		}
		acc, err := strconv.ParseInt(parts[2], 10, 32)
		if err != nil {
			return err
		}
		exp, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return err
		}
		if exp != jp.Exp {
			return errors.New("invalid timestamp")
		}
		tm := time.UnixMilli(exp)
		if tm.UTC().Before(time.Now().UTC()) {
			return errors.New("token timeout")
		}
		session.SystemId = sysId
		session.Stub = stub
		session.AccessControl = int32(acc)
		session.Successful = true
		return nil
	})
	if err != nil {
		return session, err
	}
	return session, nil
}
