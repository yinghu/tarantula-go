package main

import (
	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
)

type CacheService struct {
	bootstrap.AppManager
}

func (s *CacheService) Config() string {
	return "/etc/tarantula/cache-conf.json"
}

func (s *CacheService) Start(f conf.Env, c core.Cluster) error {
	//s.ItemUpdater = s
	s.AppManager.Start(f, c)
	//s.createSchema()
	return nil
}
