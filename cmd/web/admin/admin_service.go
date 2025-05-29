package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/event"

	"gameclustering.com/internal/bootstrap"
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

	//http.Handle("/admin",logging(&AdminLogin{AdminService: s}))
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

func logging(s bootstrap.TarantulaApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			dur := time.Since(start)
			ms := metrics.ReqMetrics{Path: r.URL.Path, ReqTimed: dur.Milliseconds(), Node:s.Cluster().Local().Name}
			s.WebRequest(ms)
		}()
		s.ServeHTTP(w,r)
	}
}
