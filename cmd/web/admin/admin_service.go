package main

import (
	"fmt"
	"log"
	"net/http"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/event"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/metrics"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type AdminService struct {
	cls  cluster.Cluster
	sql  persistence.Postgresql
	Metr metrics.MetricsService
	Tkn  core.Jwt
}

func (s *AdminService) Config() string {
	return "/etc/tarantula/admin-conf.json"
}

func (s *AdminService) Start(f conf.Env, c cluster.Cluster) error {
	s.cls = c
	sql := persistence.Postgresql{Url: f.Pgs.DatabaseURL}
	err := sql.Create()
	if err != nil {
		return err
	}
	s.sql = sql
	ms := persistence.MetricsDB{Sql: &sql}
	s.Metr = &ms

	tkn := util.JwtHMac{Alg: "SHS256"}
	tkn.HMac()
	s.Tkn = &tkn

	hash, err := util.HashPassword("password")
	if err != nil {
		return err
	}
	err = s.SaveLogin(&event.Login{Name: "root", Hash: hash})
	if err != nil {
		fmt.Printf("Root already existed %s\n", err.Error())
	}
	http.Handle("/", bootstrap.Logging(&AdminWeb{AdminService: s}))

	http.Handle("/admin/login", bootstrap.Logging(&AdminLogin{AdminService: s}))
	log.Fatal(http.ListenAndServe(f.HttpEndpoint, nil))
	return nil
}

func (s *AdminService) Metrics() metrics.MetricsService {
	return s.Metr
}
func (s *AdminService) Cluster() cluster.Cluster {
	return s.cls
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
