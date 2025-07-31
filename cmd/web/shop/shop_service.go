package main

import (
	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
)

type ShopService struct {
	bootstrap.AppManager
}

func (s *ShopService) Config() string {
	return "/etc/tarantula/shop-conf.json"
}

func (s *ShopService) Start(f conf.Env, c core.Cluster) error {
	s.AppManager.Start(f, c)
	//s.createSchema()
	return nil
}
