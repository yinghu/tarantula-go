package main

import (
	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
)

type TournamentService struct {
	bootstrap.AppManager
}

func (s *TournamentService) Config() string {
	return "/etc/tarantula/tournament-conf.json"
}

func (s *TournamentService) Start(f conf.Env, c cluster.Cluster) error {
	s.AppManager.Start(f, c)
	s.createSchema()
	return nil
}
