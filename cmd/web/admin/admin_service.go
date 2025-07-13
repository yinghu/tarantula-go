package main

import (
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type AdminService struct {
	bootstrap.AppManager
}

func (s *AdminService) Config() string {
	return "/etc/tarantula/admin-conf.json"
}

func (s *AdminService) Start(f conf.Env, c core.Cluster) error {
	s.AppManager.Start(f, c)
	err := s.createSchema()
	if err != nil {
		return err
	}
	hash, err := s.Authenticator().HashPassword("password")
	if err != nil {
		return err
	}
	err = s.SaveLogin(&event.Login{Name: "root", Hash: hash, AccessControl: bootstrap.SUDO_ACCESS_CONTROL})
	if err != nil {
		fmt.Printf("Root already existed %s\n", err.Error())
	}

	http.Handle("/admin/webprotected/{name}", bootstrap.Logging(&AdminWebProtected{AdminService: s}))
	http.Handle("/admin/web/{name}", bootstrap.Logging(&AdminWebIndex{AdminService: s}))
	http.Handle("/admin/loadenum/{name}", bootstrap.Logging(&AdminLoadEnum{AdminService: s}))
	http.Handle("/admin/saveenum", bootstrap.Logging(&AdminSaveEnum{AdminService: s}))
	http.Handle("/admin/loadcategory/{id}/{name}", bootstrap.Logging(&AdminLoadCategory{AdminService: s}))
	http.Handle("/admin/loadcategories/{scope}", bootstrap.Logging(&AdminLoadCategories{AdminService: s}))
	http.Handle("/admin/savecategory", bootstrap.Logging(&AdminSaveCategory{AdminService: s}))
	http.Handle("/admin/loadconfig/{id}", bootstrap.Logging(&AdminLoadConfig{AdminService: s}))
	http.Handle("/admin/loadconfigs/{name}/{limit}", bootstrap.Logging(&AdminLoadConfigs{AdminService: s}))
	http.Handle("/admin/saveconfig", bootstrap.Logging(&AdminSaveConfig{AdminService: s}))
	http.Handle("/admin/configapp/{app}", bootstrap.Logging(&AdminConfigApp{AdminService: s}))
	http.Handle("/admin/resetkey", bootstrap.Logging(&AdminResetKey{AdminService: s}))
	http.Handle("/admin/getnode/{group}/{name}", bootstrap.Logging(&AdminGetNode{AdminService: s}))
	http.Handle("/admin/addlogin", bootstrap.Logging(&SudoAddLogin{AdminService: s}))
	http.Handle("/admin/confignode", bootstrap.Logging(&SudoConfigNode{AdminService: s}))
	http.Handle("/admin/password", bootstrap.Logging(&AdminChangePwd{AdminService: s}))
	http.Handle("/admin/login", bootstrap.Logging(&AdminLogin{AdminService: s}))
	fmt.Printf("Admin service started %s\n", f.HttpBinding)
	return nil
}
