package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type Service struct {
	Cluster *cluster.Etc
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
	fmt.Printf("Presence service started\n")
	return nil
}
func (s *Service) Shutdown() {
	s.Sql.Close()
	fmt.Printf("Presence service shut down\n")
}

func (s *Service) Register(login *Login) {
	id, _ := s.Sfk.Id()
	login.SystemId = id
	hash, _ := util.Hash(login.Hash)
	login.Hash = hash
	err := s.SaveLogin(login)
	if err != nil {
		login.Listener <- event.Chunk{Remaining: false, Data: errorMessage(err.Error(), DB_OP_ERR_CODE)}
		return
	}
	login.Listener <- event.Chunk{Remaining: false, Data: successMessage("registered")}
}

func (s *Service) VerifyToken(token string, listener chan event.Chunk) {
	err := s.Tkn.Verify(token, func(h *util.JwtHeader, p *util.JwtPayload) error {
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

func (s *Service) Login(login *Login) {
	pwd := login.Hash
	err := s.LoadLogin(login)
	if err != nil {
		login.Listener <- event.Chunk{Remaining: false, Data: errorMessage(err.Error(), DB_OP_ERR_CODE)}
		return
	}
	er := util.Match(pwd, login.Hash)
	if er != nil {
		login.Listener <- event.Chunk{Remaining: false, Data: errorMessage(er.Error(), WRONG_PASS_CODE)}
		return
	}
	tk, trr := s.Tkn.Token(func(h *util.JwtHeader, p *util.JwtPayload) error {
		h.Kid = "kid"
		p.Aud = "player"
		exp := time.Now().Add(time.Hour * 24).UTC()
		p.Exp = exp.UnixMilli()
		return nil
	})
	if trr != nil {
		login.Listener <- event.Chunk{Remaining: false, Data: errorMessage(trr.Error(), INVALID_TOKEN_CODE)}
		return
	}
	session := OnSession{Successful: true, SystemId: login.SystemId, Stub: login.SystemId, Token: tk, Home: s.Cluster.Local.HttpEndpoint}
	login.Listener <- event.Chunk{Remaining: false, Data: util.ToJson(session)}
}

func notsupport(listener chan event.Chunk) {
	listener <- event.Chunk{Remaining: false, Data: []byte("not supported")}
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.Header.Get("Tarantula-action")
	token := r.Header.Get("Tarantula-token")
	listener := make(chan event.Chunk)
	defer func() {
		close(listener)
		r.Body.Close()
	}()
	w.WriteHeader(http.StatusOK)
	switch action {
	case "onRegister":
		var login Login
		json.NewDecoder(r.Body).Decode(&login)
		login.EventObj.Listener = listener
		go s.Register(&login)

	case "onLogin":
		var login Login
		login.EventObj.Listener = listener
		json.NewDecoder(r.Body).Decode(&login)
		go s.Login(&login)

	case "onPassword":
		go s.VerifyToken(token, listener)

	default:
		go notsupport(listener)
	}
	for c := range listener {
		w.Write(c.Data)
		if !c.Remaining {
			break
		}
	}
}
