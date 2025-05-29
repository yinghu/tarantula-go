package main

import (
	"net/http"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/event"
)

type AdminLogin struct {
	*AdminService
}

func (s AdminLogin) Login(login *event.Login) {

}

func (s AdminLogin) Cluster() cluster.Cluster {
	return s.AdminService.Cluster
}

func (s *AdminLogin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))

}
