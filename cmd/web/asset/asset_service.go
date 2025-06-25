package main

import (
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
)

type AssetService struct {
	bootstrap.AppManager
	assetDir string
}

func (s *AssetService) Config() string {
	return "/etc/tarantula/asset-conf.json"
}

func (s *AssetService) Start(f conf.Env, c core.Cluster) error {
	s.AppManager.Start(f, c)
	s.assetDir = f.LocalDir
	err := s.createSchema()
	if err != nil {
		return nil
	}
	fmt.Printf("Asset service started %s %s\n", f.HttpBinding, f.LocalDir)
	http.Handle("/asset/upload", bootstrap.Logging(&AssetUpload{AssetService: s}))
	http.Handle("/asset/download", bootstrap.Logging(&AssetDownload{AssetService: s}))
	return nil
}
