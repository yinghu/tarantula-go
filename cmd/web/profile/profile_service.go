package main

import (
	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
)

type ProfileService struct {
	bootstrap.AppManager
}

func (s *ProfileService) Config() string {
	return "/etc/tarantula/profile-conf.json"
}

func (s *ProfileService) Start(f conf.Env, c cluster.Cluster) error {
	s.AppManager.Start(f, c)
	s.createSchema()
	return nil
}
