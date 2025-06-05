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
	return "/etc/tarantula/admin-conf.json"
}

func (s *AssetService) Start(f conf.Env, c cluster.Cluster) error {
	s.AppManager.Start(f, c)

	log.Fatal(http.ListenAndServe(f.HttpBinding, nil))
	return nil
}
