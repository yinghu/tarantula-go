package main

import (
	"encoding/json"
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
	sql     persistence.Postgresql
	Seq     core.Sequence
	Tkn     core.Jwt
	Ciph    util.Cipher
	Ds      core.DataStore
	Started bool
}

func (s *PresenceService) Create(classId int) event.Event {
	login := event.Login{}
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
	sfk := util.NewSnowflake(env.NodeId, util.EpochMillisecondsFromMidnight(2020, 1, 1))
	s.Seq = &sfk
	tkn := util.JwtHMac{Alg: "SHS256"}
	tkn.HMac()
	s.Tkn = &tkn
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
	s.sql = sql
	ds := persistence.Cache{InMemory: env.Bdg.InMemory, Path: env.Bdg.Path, Seq: s.Seq, KeySize: env.Bdg.KeySize, ValueSize: env.Bdg.ValueSize}
	err = ds.Open()
	if err != nil {
		return err
	}
	s.Ds = &ds
	s.Started = true
	fmt.Printf("Presence service started\n")
	http.Handle("/presence",logging(s))
	log.Fatal(http.ListenAndServe(env.HttpEndpoint, nil))
	return nil
}
func (s *PresenceService) Shutdown() {
	s.sql.Close()
	err := s.Ds.Close()
	if err != nil {
		fmt.Printf("Error %s\n", err.Error())
	}
	fmt.Printf("Presence service shut down\n")
}

func (s *PresenceService) Publish(e event.Event) error {
	err := s.Ds.Save(e)
	if err != nil {
		fmt.Printf("Cache Error %s\n", err.Error())
	}
	v, ok := e.(*event.Login)
	if ok {
		fmt.Printf("LOGIN %s\n", v.Name)
	}
	load := event.Login{Name: v.Name}
	err = s.Ds.Load(&load)
	if err == nil {
		fmt.Printf("LOADED %d\n", load.SystemId)
	}
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
		var login event.Login
		json.NewDecoder(r.Body).Decode(&login)
		login.EventObj.Cc = listener
		go s.Register(&login)

	case "onLogin":
		var login event.Login
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
