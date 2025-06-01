package main

import (
	
	"fmt"
	"log"
	"net/http"
	

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/metrics"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type PresenceService struct {
	cls     cluster.Cluster
	Metr    metrics.MetricsService
	sql     persistence.Postgresql
	Seq     core.Sequence
	Auth    core.Authenticator
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

func (s *PresenceService) Metrics() metrics.MetricsService {
	return s.Metr
}
func (s *PresenceService) Cluster() cluster.Cluster {
	return s.cls
}
func (s *PresenceService) Authenticator() core.Authenticator {
	return s.Auth
}

func (s *PresenceService) Start(env conf.Env, c cluster.Cluster) error {
	s.cls = c
	sfk := util.NewSnowflake(env.NodeId, util.EpochMillisecondsFromMidnight(2020, 1, 1))
	s.Seq = &sfk
	tkn := util.JwtHMac{Alg: "SHS256"}
	tkn.HMac()

	ci := util.Aes{Ksz: 32}
	err := ci.AesGcm()
	if err != nil {
		return err
	}

	s.Auth = &bootstrap.AuthManager{Tkn: &tkn, Cipher: &ci, Kid: "presence", DurHours: 24}
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
	ms := persistence.MetricsDB{Sql: &sql}
	s.Metr = &ms

	s.Started = true
	fmt.Printf("Presence service started\n")
	http.Handle("/presence/register", bootstrap.Logging(&PresenceRegister{PresenceService: s}))
	http.Handle("/presence/login", bootstrap.Logging(&PresenceLogin{PresenceService: s}))
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
	for v := range s.cls.View() {
		if v.Name != s.cls.Local().Name {
			go func() {
				pub := event.SocketPublisher{Remote: v.TcpEndpoint, BufferSize: 1024}
				pub.Publish(e)
			}()
		}
	}
	return nil
}

