package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
)

type InventoryService struct {
	bootstrap.AppManager
}

func (s *InventoryService) Config() string {
	return "/etc/tarantula/inventory-conf.json"
}

func (s *InventoryService) Start(f conf.Env, c core.Cluster) error {
	s.ItemUpdater = s
	s.AppManager.Start(f, c)
	s.createSchema()
	http.Handle("/inventory/grant", bootstrap.Logging(&InventoryGranter{InventoryService: s}))
	return nil
}
