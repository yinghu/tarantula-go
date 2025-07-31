package main

import (
	"fmt"
	"net/http"
	"os"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type AdminService struct {
	bootstrap.AppManager
	assetDir    string
	contentDir  string
	publishDir  string
	managedApps []string
}

func (s *AdminService) Config() string {
	return "/etc/tarantula/admin-conf.json"
}

func (s *AdminService) Start(f conf.Env, c core.Cluster) error {
	s.AppManager.Start(f, c)
	s.managedApps = f.ManagedApps
	s.contentDir = f.Bin
	s.assetDir = f.LocalDir + "/asset"
	os.MkdirAll(s.assetDir, 0755)
	s.publishDir = f.Bin + "/tarantula"
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
		core.AppLog.Printf("Root already existed %s\n", err.Error())
	}

	http.Handle("/admin/webprotected/{name}", bootstrap.Logging(&AdminWebProtected{AdminService: s}))
	http.Handle("/admin/web/{name}", bootstrap.Logging(&AdminWebIndex{AdminService: s}))
	//handle / context from nginx proxy
	http.Handle("/admin/{name}", bootstrap.Logging(&AdminWebIndex{AdminService: s}))

	http.Handle("/admin/enum/load/{name}", bootstrap.Logging(&EnumLoader{AdminService: s}))
	http.Handle("/admin/file/save", bootstrap.Logging(&FileSaver{AdminService: s}))
	http.Handle("/admin/enum/save", bootstrap.Logging(&EnumSaver{AdminService: s}))
	http.Handle("/admin/category/load/{id}/{name}/{to}/{target}", bootstrap.Logging(&CategoryLoader{AdminService: s}))
	http.Handle("/admin/category/save", bootstrap.Logging(&CategorySaver{AdminService: s}))
	http.Handle("/admin/config/load/{id}/{name}/{limit}", bootstrap.Logging(&ConfigLoader{AdminService: s}))
	http.Handle("/admin/config/save", bootstrap.Logging(&ConfigSaver{AdminService: s}))
	http.Handle("/admin/config/register/{opt}/{id}/{app}", bootstrap.Logging(&ConfigRegister{AdminService: s}))
	http.Handle("/admin/category/preview/{id}", bootstrap.Logging(&CategoryPreviewer{AdminService: s}))
	http.Handle("/admin/item/delete/{type}/{id}", bootstrap.Logging(&ItemDeleter{AdminService: s}))

	http.Handle("/admin/view/{id}", bootstrap.Logging(&AdminItemViewer{AdminService: s}))
	http.Handle("/admin/sync/{repo}", bootstrap.Logging(&AdminPublisher{AdminService: s}))
	http.Handle("/admin/env", bootstrap.Logging(&AdminEnv{AdminService: s}))
	http.Handle("/admin/snowflake/parse/{id}", bootstrap.Logging(&AdminParseSnowFlakeId{AdminService: s}))
	http.Handle("/admin/login/add", bootstrap.Logging(&SudoAddLogin{AdminService: s}))
	http.Handle("/admin/password", bootstrap.Logging(&AdminChangePwd{AdminService: s}))
	http.Handle("/admin/login", bootstrap.Logging(&AdminLogin{AdminService: s}))
	fmt.Printf("Admin service started %s\n", f.HttpBinding)
	return nil
}
