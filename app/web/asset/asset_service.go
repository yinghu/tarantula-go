package main

import (
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"
)

type AssetService struct {
	bootstrap.AppManager
	assetDir string
}

func (s *AssetService) Config() string {
	return "./asset-conf.json"
}

func (s *AssetService) Start(f core.Env) error {
	s.AppManager.Start(f)
	s.assetDir = fmt.Sprintf("%s/%s", f.HomeDir, f.GroupName)
	err := s.createSchema()
	if err != nil {
		return nil
	}
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
	s.Cluster().Register("update", &protocol.TccTransationListener{Reserve: func(e *protocol.Transaction) error {
		return nil 
	}, Confirm: func(e *protocol.Transaction) error {
		return nil
	}, Cancel: func(e *protocol.Transaction) error {
		return nil
	}})
	core.AppLog.Printf("Asset service started %s %s\n", f.HttpBinding, s.assetDir)
	http.Handle("/asset/upload/{name}", bootstrap.Logging(&AssetUpload{AssetService: s}))
	http.Handle("/asset/download/{name}", bootstrap.Logging(&AssetDownload{AssetService: s}))
	return nil
}
