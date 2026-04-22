package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
	"google.golang.org/protobuf/proto"
)

type PresenceService struct {
	bootstrap.AppManager
	Started     bool
	LoginReward core.Configuration
}

func (s *PresenceService) Config() string {
	return "/etc/tarantula/presence-conf.json"
}

func (s *PresenceService) Start(env core.Env) error {
	s.ItemUpdater = s
	err := s.AppManager.Start(env)
	if err != nil {
		return err
	}
	err = s.createSchema()
	if err != nil {
		return err
	}
	brn := util.GitCurBranch()
	regs, err := s.ItemService().LoadRegistrations(s.Context(), brn.Message)
	if err == nil {
		for i := range regs {
			c, err := s.ItemService().InventoryManager().Load(regs[i].ItemId)
			if err == nil {
				s.ItemListener().OnRegister(c)
			} else {
				core.AppLog.Printf("Error on load registration %s %s\n", err.Error(), brn.Message)
			}
		}
	} else {
		core.AppLog.Printf("Error on load registration %s %s\n", err.Error(), brn.Message)
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
		core.AppLog.Debug().Msgf("reserve %v", e)
		return nil //fmt.Errorf("no")
	}, Confirm: func(e *protocol.Transaction) error {
		core.AppLog.Debug().Msgf("confirm %v", e)
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

func (s *PresenceService) OnEvent(e core.Event) {
	core.AppLog.Printf("%v\n", e)
}
