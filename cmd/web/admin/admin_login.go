package main

import (
	"net/http"
	"time"

	"gameclustering.com/internal/event"
	"gameclustering.com/internal/metrics"
)

type AdminLogin struct {
	*AdminService
}

func (s AdminLogin) Login(login *event.Login) {
	
}

func (s *AdminLogin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))

}

func logging(s *AdminLogin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			dur := time.Since(start)
			ms := metrics.ReqMetrics{Path: r.URL.Path, ReqTimed: dur.Milliseconds(), Node: s.Cluster.Local().Name}
			s.Metr.WebRequest(ms)
		}()
		s.ServeHTTP(w, r)
	}
}
