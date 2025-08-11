package main

import (
	"net/http"
	"sync"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type TournamentMap struct {
	sync.RWMutex
	tournaments map[int64]any
}

type TournamentService struct {
	bootstrap.AppManager
	TournamentMap
}

func (s *TournamentService) Config() string {
	return "/etc/tarantula/tournament-conf.json"
}

func (s *TournamentService) Start(f conf.Env, c core.Cluster) error {
	s.ItemUpdater = s
	s.AppManager.Start(f, c)
	s.createSchema()
	s.tournaments = make(map[int64]any)
	http.Handle("/tournament/list", bootstrap.Logging(&TournamentList{TournamentService: s}))
	return nil
}

func (s *TournamentService) OnEvent(e event.Event) {
	core.AppLog.Printf("%v\n", e)
}
