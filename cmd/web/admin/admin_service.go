package main

import (
	"fmt"
	"log"
	"net/http"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/event"

	"gameclustering.com/internal/bootstrap"
)

type AdminService struct {
	bootstrap.AppManager
}

func (s *AdminService) Config() string {
	return "/etc/tarantula/admin-conf.json"
}

func (s *AdminService) Start(f conf.Env, c cluster.Cluster) error {
	s.AppManager.Start(f, c)
	hash, err := s.Auth.HashPassword("password")
	if err != nil {
		return err
	}
	err = s.SaveLogin(&event.Login{Name: "root", Hash: hash})
	if err != nil {
		fmt.Printf("Root already existed %s\n", err.Error())
	}
	http.Handle("/admin/confignode", bootstrap.Logging(&SudoConfigNode{AdminService: s}))
	http.Handle("/admin/password", bootstrap.Logging(&AdminChangePwd{AdminService: s}))
	http.Handle("/admin/login", bootstrap.Logging(&AdminLogin{AdminService: s}))
	fmt.Printf("Admin service started %s\n", f.HttpBinding)
	log.Fatal(http.ListenAndServe(f.HttpBinding, nil))
	return nil
}
