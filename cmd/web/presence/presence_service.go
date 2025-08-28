package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/util"
)

type PresenceService struct {
	bootstrap.AppManager
	Started     bool
	LoginReward item.Configuration
}

func (s *PresenceService) Config() string {
	return "/etc/tarantula/presence-conf.json"
}

func (s *PresenceService) Start(env conf.Env, c core.Cluster) error {
	s.ItemUpdater = s
	err := s.AppManager.Start(env, c)
	if err != nil {
		return err
	}
	err = s.createSchema()
	if err != nil {
		return err
	}
	brn := util.GitCurBranch()
	core.AppLog.Printf("Item registratuin %s %s\n", s.Context(), brn.Message)
	regs, err := s.ItemService().LoadRegistrations(s.Context(), brn.Message)
	if err == nil {
		core.AppLog.Printf("Item %d\n", len(regs))
		for i := range regs {
			core.AppLog.Printf("Item %d\n", regs[i].ItemId)
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
	core.AppLog.Printf("Presence service started %s\n", env.HttpBinding)
	http.Handle("/presence/register", bootstrap.Logging(&PresenceRegister{PresenceService: s}))
	http.Handle("/presence/login", bootstrap.Logging(&PresenceLogin{PresenceService: s}))
	http.Handle("/presence/password", bootstrap.Logging(&PresenceChangePwd{PresenceService: s}))
	return nil
}
func (s *PresenceService) Shutdown() {
	s.AppManager.Shutdown()
	core.AppLog.Printf("Presence service shut down\n")
}

func (s *PresenceService) OnEvent(e event.Event) {
	core.AppLog.Printf("%v\n", e)
}
