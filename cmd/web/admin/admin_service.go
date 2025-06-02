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
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type AdminService struct {
	bootstrap.AppManager
	sql  persistence.Postgresql
}

func (s *AdminService) Config() string {
	return "/etc/tarantula/admin-conf.json"
}

func (s *AdminService) Start(f conf.Env, c cluster.Cluster) error {
	s.Cls = c
	sql := persistence.Postgresql{Url: f.Pgs.DatabaseURL}
	err := sql.Create()
	if err != nil {
		return err
	}
	s.sql = sql
	ms := persistence.MetricsDB{Sql: &sql}
	s.Metr = &ms

	tkn := util.JwtHMac{Alg: "SHS256", Ksz: core.JWT_KEY_SIZE}
	ci := util.Aes{Ksz: core.CIPHER_KEY_SIZE}
	err = c.Atomic(func(ctx cluster.Ctx) error {
		jsk, err := ctx.Get(core.JWT_KEY_NAME)
		if err != nil {
			fmt.Println("Create new jwt key")
			nkey := util.Key(tkn.Ksz)
			ctx.Put(core.JWT_KEY_NAME, util.KeyToBase64(nkey))
			tkn.HMacFromKey(nkey)
			return nil
		}
		jk, err := util.KeyFromBase64(jsk)
		if err != nil {
			return err
		}
		tkn.HMacFromKey(jk)
		return nil
	})
	if err != nil {
		return err
	}
	err = c.Atomic(func(ctx cluster.Ctx) error {
		csk, err := ctx.Get(core.CIPHER_KEY_NAME)
		if err != nil {
			fmt.Println("Create new cipher key")
			ckey := util.Key(ci.Ksz)
			ctx.Put(core.CIPHER_KEY_NAME, util.KeyToBase64(ckey))
			ci.AesGcmFromKey(ckey)
		}
		ck, err := util.KeyFromBase64(csk)
		if err != nil {
			return err
		}
		ci.AesGcmFromKey(ck)
		return nil
	})
	if err != nil {
		return err
	}

	s.Auth = &bootstrap.AuthManager{Tkn: &tkn, Cipher: &ci, Kid: "admin"}

	hash, err := s.Auth.HashPassword("password")
	if err != nil {
		return err
	}
	err = s.SaveLogin(&event.Login{Name: "root", Hash: hash})
	if err != nil {
		fmt.Printf("Root already existed %s\n", err.Error())
	}
	http.Handle("/admin/password", bootstrap.Logging(&AdminChangePwd{AdminService: s}))
	http.Handle("/admin/login", bootstrap.Logging(&AdminLogin{AdminService: s}))
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
