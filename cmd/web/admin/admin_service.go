package main

import (
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/event"
)

type AdminService struct {
	bootstrap.AppManager
}

func (s *AdminService) Config() string {
	return "/etc/tarantula/admin-conf.json"
}

func (s *AdminService) Start(f conf.Env, c cluster.Cluster) error {
	s.AppManager.Start(f, c)
	hash, err := s.Authenticator().HashPassword("password")
	if err != nil {
		return err
	}
	err = s.SaveLogin(&event.Login{Name: "root", Hash: hash, AccessControl: bootstrap.SUDO_ACCESS_CONTROL})
	if err != nil {
		fmt.Printf("Root already existed %s\n", err.Error())
	}
	configApp := AdminConfigApp{AdminService: s}
	err = configApp.start()
	if err != nil {
		return err
	}
	http.Handle("/admin/configapp/{app}", bootstrap.Logging(&configApp))
	http.Handle("/admin/resetkey", bootstrap.Logging(&AdminResetKey{AdminService: s}))
	http.Handle("/admin/getnode/{group}/{name}", bootstrap.Logging(&AdminGetNode{AdminService: s}))
	http.Handle("/admin/addlogin", bootstrap.Logging(&SudoAddLogin{AdminService: s}))
	http.Handle("/admin/confignode", bootstrap.Logging(&SudoConfigNode{AdminService: s}))
	http.Handle("/admin/password", bootstrap.Logging(&AdminChangePwd{AdminService: s}))
	http.Handle("/admin/login", bootstrap.Logging(&AdminLogin{AdminService: s}))
	fmt.Printf("Admin service started %s\n", f.HttpBinding)
	return nil
}
