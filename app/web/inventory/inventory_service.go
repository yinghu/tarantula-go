package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type InventoryService struct {
	bootstrap.AppManager
}

func (s *InventoryService) Config() string {
	return "/etc/tarantula/inventory-conf.json"
}

func (s *InventoryService) Start(f conf.Env, c core.Cluster, p event.Pusher) error {
	s.ItemUpdater = s
	s.AppManager.Start(f, c, p)
	s.createSchema()
	http.Handle("/inventory/grant", bootstrap.Logging(&InventoryGranter{InventoryService: s}))
	http.Handle("/inventory/load", bootstrap.Logging(&InventoryLoader{InventoryService: s}))
	return nil
}
