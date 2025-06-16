package main

import (
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type PresenceService struct {
	bootstrap.AppManager
	Seq     core.Sequence
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
		return
	}
}

func (s *PresenceService) Config() string {
	return "/etc/tarantula/presence-conf.json"
}

func (s *PresenceService) Start(env conf.Env, c cluster.Cluster) error {
	err := s.AppManager.Start(env, c)
	if err != nil {
		return err
	}
	sfk := util.NewSnowflake(env.NodeId, util.EpochMillisecondsFromMidnight(2020, 1, 1))
	s.Seq = &sfk
	path := env.LocalDir + "/store"
	ds := persistence.Cache{InMemory: env.Bdg.InMemory, Path: path, Seq: s.Seq}
	err = ds.Open()
	if err != nil {
		return err
	}
	s.Ds = &ds
	s.Started = true
	fmt.Printf("Presence service started %s\n", env.HttpBinding)
	http.Handle("/presence/admin", bootstrap.Logging(&PresenceAdmin{PresenceService: s}))
	http.Handle("/presence/register", bootstrap.Logging(&PresenceRegister{PresenceService: s}))
	http.Handle("/presence/login", bootstrap.Logging(&PresenceLogin{PresenceService: s}))
	http.Handle("/presence/password", bootstrap.Logging(&PresenceChangePwd{PresenceService: s}))
	return nil
}
func (s *PresenceService) Shutdown() {
	s.AppManager.Shutdown()
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
	for v := range s.Cluster().View() {
		if v.Name != s.Cluster().Local().Name {
			go func() {
				pub := event.SocketPublisher{Remote: v.TcpEndpoint}
				pub.Publish(e)
			}()
		}
	}
	return nil
}
