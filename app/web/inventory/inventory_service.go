package main

import (
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

type InventoryService struct {
	bootstrap.AppManager
}

func (s *InventoryService) Config() string {
	return "/etc/tarantula/inventory-conf.json"
}

func (s *InventoryService) Start(f core.Env) error {
	s.ItemUpdater = s
	s.AppManager.Start(f)
	s.createSchema()
	s.Cluster().Subscribe("message", &protocol.MessageEventListener{Callback: func(e *protocol.MessageEvent) error {
		core.AppLog.Debug().Msgf("REV : %v", e)
		return nil
	}})
	s.Cluster().Register("grant", &protocol.TccTransationListener{Reserve: func(e *protocol.Transaction) error {
		core.AppLog.Debug().Msgf("reserve resource %v", e)
		return fmt.Errorf("no resource")
	}, Confirm: func(e *protocol.Transaction) error {
		core.AppLog.Debug().Msgf("confirm resource %v", e)
		return nil
	}, Cancel: func(e *protocol.Transaction) error {
		core.AppLog.Debug().Msgf("cancel resource %v", e)
		return nil
	}})
	http.Handle("/inventory/grant", bootstrap.Logging(&InventoryGranter{InventoryService: s}))
	http.Handle("/inventory/load", bootstrap.Logging(&InventoryLoader{InventoryService: s}))
	http.Handle("/inventory/cluster/update", bootstrap.Logging(&InventoryClusterUpdate{InventoryService: s}))
	return nil
}
