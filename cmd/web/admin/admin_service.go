package main

import (
	"fmt"
	"log"
	"net/http"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/event"

	//"gameclustering.com/internal/core"
	//"gameclustering.com/internal/event"
	"gameclustering.com/internal/metrics"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type AdminService struct {
	Cluster cluster.Cluster
	sql     persistence.Postgresql
	Metr      metrics.MetricsService
}

func (s *AdminService) Config() string {
	return "/etc/tarantula/admin-conf.json"
}

func (s *AdminService) Start(f conf.Env, c cluster.Cluster) error {
	s.Cluster = c
	sql := persistence.Postgresql{Url: f.Pgs.DatabaseURL}
	err := sql.Create()
	if err != nil {
		return err
	}
	s.sql = sql
	ms := persistence.MetricsDB{Sql: &sql}
	s.Metr = &ms
	hash, err := util.HashPassword("password")
	if err != nil {
		return err
	}
	err = s.SaveLogin(&event.Login{Name: "root", Hash: hash})
	if err != nil {
		fmt.Printf("Root already existed %s\n", err.Error())
	}
	http.HandleFunc("/", handleWeb)

	http.Handle("/admin",logging(&AdminLogin{AdminService: s}))
	log.Fatal(http.ListenAndServe(f.HttpEndpoint, nil))
	return nil
}

func (s *AdminService) Shutdown() {
	s.sql.Close()
	fmt.Printf("Admin service shut down\n")
}

func (s *AdminService) Create(classId int) event.Event {
	return &event.Login{}
}

func (s *AdminService) OnEvent(e event.Event) {

}
