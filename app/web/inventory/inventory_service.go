package main

import (
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"
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
	s.Cluster().Subscribe(event.MESSAGE_TOPIC_NAME, &protocol.TopicEventListener{C: func() proto.Message {
		return &protocol.MessageEvent{}
	}, M: func(m proto.Message) {
		ro, ok := m.(*protocol.MessageEvent)
		if ok {
			core.AppLog.Debug().Msgf("MESSAGE event %s %s", ro.Message, ro.Source)
		} else {
			core.AppLog.Debug().Msg("wrong type")
		}
	}})
	s.Cluster().Register("grant", &protocol.TccTransationListener{Reserve: func(e *protocol.Transaction) error {
		mf := persistence.NewCommodityObjectFactory()
		obj, err := mf.Message(e.Object)
		if err != nil {
			core.AppLog.Warn().Msgf("wrong decode format %s", err)
			return err
		}
		m, ok := obj.(*protocol.Commodity)
		if !ok {
			return fmt.Errorf("wrong data passed from task")
		}
		core.AppLog.Debug().Msgf("%v", m)
		return nil
	}, Confirm: func(e *protocol.Transaction) error {

		return nil
	}, Cancel: func(e *protocol.Transaction) error {

		return nil
	}})
	http.Handle("/inventory/grant", bootstrap.Logging(&InventoryGranter{InventoryService: s}))
	http.Handle("/inventory/load", bootstrap.Logging(&InventoryLoader{InventoryService: s}))
	http.Handle("/inventory/cluster/update", bootstrap.Logging(&InventoryClusterUpdate{InventoryService: s}))
	return nil
}
