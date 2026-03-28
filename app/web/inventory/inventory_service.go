package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
)

type InventoryService struct {
	bootstrap.AppManager
}

func (s *InventoryService) Config() string {
	return "/etc/tarantula/inventory-conf.json"
}

func (s *InventoryService) Start(f core.Env, p core.Pusher) error {
	s.ItemUpdater = s
	s.AppManager.Start(f, p)
	s.createSchema()
	s.Cluster().Subscribe("login", s.Event())
	http.Handle("/inventory/grant", bootstrap.Logging(&InventoryGranter{InventoryService: s}))
	http.Handle("/inventory/load", bootstrap.Logging(&InventoryLoader{InventoryService: s}))
	http.Handle("/inventory/cluster/update", bootstrap.Logging(&InventoryClusterUpdate{InventoryService: s}))
	return nil
}
