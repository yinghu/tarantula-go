package main

import (
	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
)

type ProfileService struct {
	bootstrap.AppManager
}

func (s *ProfileService) Config() string {
	return "/etc/tarantula/profile-conf.json"
}

func (s *ProfileService) Start(f conf.Env, c core.Cluster) error {
	s.AppManager.Start(f, c)
	s.createSchema()
	return nil
}
