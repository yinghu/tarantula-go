package main

import (
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

type AssetService struct {
	bootstrap.AppManager
	assetDir string
}

func (s *AssetService) Config() string {
	return "/etc/tarantula/asset-conf.json"
}

func (s *AssetService) Start(f core.Env) error {
	s.ItemUpdater = s
	s.AppManager.Start(f)
	s.assetDir = fmt.Sprintf("%s/%s", f.HomeDir, f.GroupName)
	err := s.createSchema()
	if err != nil {
		return nil
	}
	s.Cluster().Subscribe("message", &protocol.MessageEventListener{Callback: func(e *protocol.MessageEvent) error {
		core.AppLog.Debug().Msgf("REV : %v", e)
		return nil
	}})
	s.Cluster().Register("update", &protocol.TccTransationListener{Reserve: func(e *protocol.Transaction) error {
		core.AppLog.Debug().Msgf("reserve update %v", e)
		return nil
	}, Confirm: func(e *protocol.Transaction) error {
		core.AppLog.Debug().Msgf("update %v", e)
		return nil
	}, Cancel: func(e *protocol.Transaction) error {
		return nil
	}})
	core.AppLog.Printf("Asset service started %s %s\n", f.HttpBinding, s.assetDir)
	http.Handle("/asset/upload/{name}", bootstrap.Logging(&AssetUpload{AssetService: s}))
	http.Handle("/asset/download/{name}", bootstrap.Logging(&AssetDownload{AssetService: s}))
	http.Handle("/asset/cluster/create", bootstrap.Logging(&AssetClusterCreate{AssetService: s}))

	return nil
}
