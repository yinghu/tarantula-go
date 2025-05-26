package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/event"

	//"gameclustering.com/internal/core"
	//"gameclustering.com/internal/event"
	//"gameclustering.com/internal/metrics"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type AdminService struct {
	Cluster cluster.Cluster
	sql     persistence.Postgresql
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
	hash, err := util.HashPassword("password")
	if err != nil {
		return err
	}
	err = s.SaveLogin(&event.Login{Name: "root", Hash: hash})
	if err != nil {
		fmt.Printf("Root already existed %s\n", err.Error())
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" || path == "/index.html" {
			f, err := os.Open("web/index.html")
			if err != nil {
				w.Write([]byte("no page"))
				return
			}
			defer f.Close()
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			io.Copy(w, f)
			return
		}
		ext := filepath.Ext(path)
		fmt.Printf("Ext %s , %s\n", ext, "web"+path)
		if ext == ".css" {

			f, err := os.Open("web" + path)
			if err != nil {
				w.Write([]byte("no page"))
				return
			}
			defer f.Close()
			w.Header().Set("Content-Type", "text/css")
			w.WriteHeader(http.StatusOK)
			io.Copy(w, f)
			return
		}
		if ext == ".js" {
			f, err := os.Open("web" + path)
			if err != nil {
				w.Write([]byte("no page"))
				return
			}
			defer f.Close()
			w.Header().Set("Content-Type", "text/javascript")
			w.WriteHeader(http.StatusOK)
			io.Copy(w, f)
			return
		}
		if ext == ".ico" {
			f, err := os.Open("web" + path)
			if err != nil {
				w.Write([]byte("no page"))
				return
			}
			defer f.Close()
			w.Header().Set("Content-Type", "text/css")
			w.WriteHeader(http.StatusOK)
			io.Copy(w, f)
			return
		}
		w.Write([]byte("no page" + ext))
	})

	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
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
