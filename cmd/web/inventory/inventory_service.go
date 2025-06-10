package main

import (
	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
)

type InventoryService struct {
	bootstrap.AppManager
}

func (s *InventoryService) Config() string {
	return "/etc/tarantula/asset-conf.json"
}

func (s *InventoryService) Start(f conf.Env, c cluster.Cluster) error {
	s.AppManager.Start(f, c)
	return nil
}
