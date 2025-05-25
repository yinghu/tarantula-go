package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/metrics"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type PresenceService struct {
	Cluster cluster.Cluster
	Sql     persistence.Postgresql
	Sfk     util.Snowflake
	Tkn     util.JwtHMac
	Ciph    util.Cipher
	Ds      core.DataStore
	Started bool
}

func (s *PresenceService) Create(classId int) event.Event {
	login := Login{}
	login.Cb = s
	return &login
}

func (s *PresenceService) OnEvent(e event.Event) {
	err := s.Ds.Save(e)
	if err != nil {
		fmt.Printf("No save %s\n", err.Error())
	}
}

func (s *PresenceService) Config() string {
	return "/etc/tarantula/presence-conf.json"
}

func (s *PresenceService) Start(env conf.Env, c cluster.Cluster) error {
	s.Cluster = c
	s.Sfk = util.NewSnowflake(env.NodeId, util.EpochMillisecondsFromMidnight(2020, 1, 1))
	s.Tkn = util.JwtHMac{Alg: "SHS256"}
	s.Tkn.HMac()
	ci := util.Cipher{Ksz: 32}
	err := ci.AesGcm()
	if err != nil {
		return err
	}
	s.Ciph = ci
	sql := persistence.Postgresql{Url: env.Pgs.DatabaseURL}
	err = sql.Create()
	if err != nil {
		return err
	}
	s.Sql = sql
	ds := persistence.Cache{InMemory: env.Bdg.InMemory, Path: env.Bdg.Path, Sfk: &s.Sfk, KeySize: env.Bdg.KeySize, ValueSize: env.Bdg.ValueSize}
	err = ds.Open()
	if err != nil {
		return err
	}
	s.Ds = &ds
	s.Started = true
	fmt.Printf("Presence service started\n")
	http.Handle("/presence", http.HandlerFunc(logging(s)))
	log.Fatal(http.ListenAndServe(env.HttpEndpoint, nil))
	return nil
}
func (s *PresenceService) Shutdown() {
	s.Sql.Close()
	s.Ds.Close()
	fmt.Printf("Presence service shut down\n")
}

func (s *PresenceService) Publish(e event.Event) error {
	for v := range s.Cluster.View() {
		if v.Name != s.Cluster.Local().Name {
			go func() {
				pub := event.SocketPublisher{Remote: v.TcpEndpoint, BufferSize: 1024}
				pub.Publish(e)
			}()
		}
	}
	return nil
}

func (s *PresenceService) Register(login *Login) {
	id, _ := s.Sfk.Id()
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

func (s *PresenceService) Login(login *Login) {
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
	tk, err := s.Tkn.Token(func(h *core.JwtHeader, p *core.JwtPayload) error {
		h.Kid = "kid"
		p.Aud = "player"
		exp := time.Now().Add(time.Hour * 24).UTC()
		p.Exp = exp.UnixMilli()
		return nil
	})
	if err != nil {
		login.Cc <- event.Chunk{Remaining: false, Data: errorMessage(err.Error(), INVALID_TOKEN_CODE)}
		return
	}
	session := OnSession{Successful: true, SystemId: login.SystemId, Stub: login.SystemId, Token: tk, Home: s.Cluster.Local().HttpEndpoint}
	login.Cc <- event.Chunk{Remaining: false, Data: util.ToJson(session)}
}

func notsupport(listener chan event.Chunk) {
	listener <- event.Chunk{Remaining: false, Data: []byte("not supported")}
}

func (s *PresenceService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		login.EventObj.Cc = listener
		go s.Register(&login)

	case "onLogin":
		var login Login
		login.EventObj.Cc = listener
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

func logging(s *PresenceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		action := r.Header.Get("Tarantula-action")
		defer func() {
			dur := time.Since(start)
			ms := metrics.ReqMetrics{Path: r.URL.Path + "/" + action, ReqTimed: dur.Milliseconds(), Node: s.Cluster.Local().Name}
			s.SaveMetrics(&ms)
		}()
		s.ServeHTTP(w, r)
	}
}
