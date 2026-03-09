package main

import (
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
)

type AssetService struct {
	bootstrap.AppManager
	assetDir string
}

func (s *AssetService) Config() string {
	return "/etc/tarantula/asset-conf.json"
}

func (s *AssetService) Start(f core.Env, p core.Pusher) error {
	s.ItemUpdater = s
	s.AppManager.Start(f, p)
	s.assetDir = fmt.Sprintf("%s/%s", f.HomeDir, f.GroupName)
	err := s.createSchema()
	if err != nil {
		return nil
	}
	core.AppLog.Printf("Asset service started %s %s\n", f.HttpBinding, s.assetDir)
	http.Handle("/asset/upload/{name}", bootstrap.Logging(&AssetUpload{AssetService: s}))
	http.Handle("/asset/download/{name}", bootstrap.Logging(&AssetDownload{AssetService: s}))
	return nil
}
