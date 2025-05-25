package main

import (

	"log"
	"net/http"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/event"
	//"gameclustering.com/internal/core"
	//"gameclustering.com/internal/event"
	//"gameclustering.com/internal/metrics"
	//"gameclustering.com/internal/persistence"
	//"gameclustering.com/internal/util"
)

type AdminService struct {
	Cluster cluster.Cluster
}

func (s *AdminService) Config() string {
	return "/etc/tarantula/presence-conf.json"
}

func (s *AdminService) Start(f conf.Env, c cluster.Cluster) error {
	s.Cluster = c
	//http.Handle("/admin", http.HandlerFunc(logging(s)))
	log.Fatal(http.ListenAndServe(f.HttpEndpoint, nil))
	return nil
}

func (s *AdminService) Shutdown() {

}

func (s *AdminService) Create(classId int) event.Event {
	return nil
}

func (s *AdminService) OnEvent(e event.Event) {

}
