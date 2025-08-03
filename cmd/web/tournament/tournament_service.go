package main

import (
	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type TournamentService struct {
	bootstrap.AppManager
}

func (s *TournamentService) Config() string {
	return "/etc/tarantula/tournament-conf.json"
}

func (s *TournamentService) Start(f conf.Env, c core.Cluster) error {
	s.ItemUpdater = s
	s.AppManager.Start(f, c)
	s.createSchema()
	return nil
}

func (s *TournamentService) OnEvent(e event.Event) {
	core.AppLog.Printf("%v\n", e)
}
