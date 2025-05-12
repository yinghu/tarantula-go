
package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type Chunk struct {
	Remaining bool
	Data      []byte
}

type Login struct {
	Name        string `json:"login"`
	Hash        string `json:"password"`
	ReferenceId int32  `json:"referenceId"`
	SystemId    int64
}

type Service struct {
	Sql     persistence.Postgresql
	Sfk     util.Snowflake
	Tkn     util.Jwt
	Ciph    util.Cipher
	Started bool
}

func (s *Service) Start(env conf.Env) error {
	s.Sfk = util.NewSnowflake(env.NodeId, util.EpochMillisecondsFromMidnight(2020, 1, 1))
	s.Tkn = util.Jwt{Alg: "SHS256"}
	s.Tkn.HMac()
	ci := util.Cipher{Ksz: 32}
	er := ci.AesGcm()
	if er != nil {
		return er
	}
	s.Ciph = ci
	sql := persistence.Postgresql{Url: env.DatabaseURL}
	err := sql.Create()
	if err != nil {
		return err
	}
	s.Sql = sql
	s.Started = true
	fmt.Printf("Auth service started\n")
	return nil
}
func (s *Service) Shutdown() {
	s.Sql.Close()
	fmt.Printf("Auth service shut down\n")
}

func (s *Service) Register(login *Login) error {
	id, _ := s.Sfk.Id()
	login.SystemId = id
	hash, _ := util.Hash(login.Hash)
	login.Hash = hash
	return s.SaveLogin(login)
}

func (s *Service) VerifyToken(token string) error {
	return s.Tkn.Verify(token, func(h *util.JwtHeader, p *util.JwtPayload) error {
		t := time.UnixMilli(p.Exp).UTC()
		if t.Before(time.Now().UTC()) {
			return errors.New("token expired")
		}
		return nil
	})
}

func (s *Service) Login(login *Login) (string, error) {
	pwd := login.Hash
	err := s.LoadLogin(login)
	if err != nil {
		return "", err
	}
	//fmt.Printf("Hash %s >> %d\n", login.Hash, login.SystemId)
	er := util.Match(pwd, login.Hash)
	if er != nil {
		return "", er
	}
	tk, trr := s.Tkn.Token(func(h *util.JwtHeader, p *util.JwtPayload) error {
		h.Kid = "kid"
		p.Aud = "player"
		exp := time.Now().Add(time.Hour * 24).UTC()
		p.Exp = exp.UnixMilli()
		return nil
	})
	if trr != nil {
		return "", trr
	}
	return tk, nil
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.Header.Get("Tarantula-action")
	token := r.Header.Get("Tarantula-token")
	defer func() {
		r.Body.Close()
	}()
	w.WriteHeader(http.StatusOK)
	switch action {
	case "onRegister":
		var login Login
		json.NewDecoder(r.Body).Decode(&login)
		err := s.Register(&login)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte("success"))
		}
	case "onLogin":
		var login Login
		json.NewDecoder(r.Body).Decode(&login)
		tk, err := s.Login(&login)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte(tk))
		}
	case "onPassword":
		err := s.VerifyToken(token)
		if err != nil {
			w.Write([]byte("bad token"))
		} else {
			w.Write([]byte(token))
		}
	default:
		w.Write([]byte("not supported"))
	}
}
