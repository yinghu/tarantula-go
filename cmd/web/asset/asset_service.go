package main

import (
	"fmt"
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
	err := s.createSchema()
	if err != nil {
		return nil
	}
	fmt.Printf("Asset service started %s %s\n", f.HttpBinding, f.LocalDir)
	http.Handle("/asset/upload", bootstrap.Logging(&AssetUpload{AssetService: s}))
	http.Handle("/asset/download", bootstrap.Logging(&AssetDownload{AssetService: s}))
	log.Fatal(http.ListenAndServe(f.HttpBinding, nil))
	return nil
}
