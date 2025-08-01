package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
)

type CacheService struct {
	bootstrap.AppManager
	Ds core.DataStore
}

func (s *CacheService) Config() string {
	return "/etc/tarantula/cache-conf.json"
}

func (s *CacheService) Start(env conf.Env, c core.Cluster) error {
	s.AppManager.Start(env, c)
	path := env.LocalDir + "/store"
	ds := persistence.Cache{InMemory: env.Bdg.InMemory, Path: path, Seq: s.Sequence()}
	err := ds.Open()
	if err != nil {
		return err
	}
	s.Ds = &ds
	core.AppLog.Printf("Cache service started %s %s\n", env.HttpBinding, env.LocalDir)
	http.Handle("/cache/set/{key}", bootstrap.Logging(&CacheSetter{CacheService: s}))
	http.Handle("/cache/get/{key}", bootstrap.Logging(&CacheGetter{CacheService: s}))
	return nil
}
