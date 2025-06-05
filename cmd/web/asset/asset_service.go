package main

import (
	"log"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
)

type AssetService struct {
	bootstrap.AppManager
}

func (s *AssetService) Config() string {
	return "/etc/tarantula/asset-conf.json"
}

func (s *AssetService) Start(f conf.Env, c cluster.Cluster) error {
	s.AppManager.Start(f, c)
	http.Handle("asset/upload", bootstrap.Logging(&AssetUpload{AssetService: s}))
	http.Handle("asset/download", bootstrap.Logging(&AssetDownload{AssetService: s}))
	log.Fatal(http.ListenAndServe(f.HttpBinding, nil))
	return nil
}
