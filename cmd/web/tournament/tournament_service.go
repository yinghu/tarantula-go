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
	tournaments map[int64]Tournament
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
	s.Bsl = s
	s.AppManager.Start(f, c)
	s.createSchema()
	s.tournaments = make(map[int64]Tournament)
	ids, err := s.loadSchedule()
	if err == nil {
		for i := range ids {
			c, err := s.ItemService().Loader().Load(ids[i])
			if err == nil {
				s.ItemListener().OnRegister(c)
			}
		}
	}
	http.Handle("/tournament/list", bootstrap.Logging(&TournamentList{TournamentService: s}))
	http.Handle("/tournament/join", bootstrap.Logging(&TournamentJoin{TournamentService: s}))
	http.Handle("/tournament/score", bootstrap.Logging(&TournamentScore{TournamentService: s}))
	http.Handle("/tournament/raceboard", bootstrap.Logging(&TournamentRaceBoard{TournamentService: s}))
	return nil
}

func (s *TournamentService) OnEvent(e event.Event) {
	te, isTe := e.(*event.TournamentEvent)
	if !isTe {
		return
	}
	tmnt := s.tournaments[te.TournamentId]
	tmnt.OnBoard(*te)
}

func (s *TournamentService) NodeStarted(n core.Node) {
	core.AppLog.Printf("Node started %s %s\n", n.Name, s.Cluster().Local().Name)
	if n.Name != s.Cluster().Local().Name {
		return
	}
	for _, t := range s.tournaments {
		t.Start()
	}
}
