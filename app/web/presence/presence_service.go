package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"
)

type PresenceService struct {
	bootstrap.AppManager
	Started bool
}

func (s *PresenceService) Config() string {
	return "./presence-conf.json"
}

func (s *PresenceService) Start(env core.Env) error {
	err := s.AppManager.Start(env)
	if err != nil {
		return err
	}
	err = s.createSchema()
	if err != nil {
		return err
	}

	s.Started = true
	s.Cluster().Subscribe("register", &protocol.TopicEventListener{C: func() proto.Message {
		return &protocol.RegisterEvent{}
	}, M: func(m proto.Message) {
		ro, ok := m.(*protocol.RegisterEvent)
		if ok {
			core.AppLog.Debug().Msgf("register event %s %s", ro.Name, ro.Source)
		} else {
			core.AppLog.Debug().Msg("wrong type")
		}
	}})
	s.Cluster().Register("register", &protocol.TccTransationListener{Reserve: func(e *protocol.Transaction) error {
		return nil
	}, Confirm: func(e *protocol.Transaction) error {
		return nil
	}, Cancel: func(e *protocol.Transaction) error {
		return nil
	}})
	core.AppLog.Printf("Presence service started %s\n", env.HttpBinding)
	http.Handle("/presence/register", bootstrap.Logging(&PresenceRegister{PresenceService: s}))
	http.Handle("/presence/login", bootstrap.Logging(&PresenceLogin{PresenceService: s}))
	http.Handle("/presence/password", bootstrap.Logging(&PresenceChangePwd{PresenceService: s}))
	http.Handle("/presence/cluster/get/{key}", bootstrap.Logging(&PresenceClusterGet{PresenceService: s}))
	return nil
}
func (s *PresenceService) Shutdown() {
	s.AppManager.Shutdown()
	core.AppLog.Printf("Presence service shut down\n")
}
