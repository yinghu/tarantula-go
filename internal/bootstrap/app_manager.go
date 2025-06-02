package bootstrap

import (
	"fmt"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/metrics"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type AppManager struct {
	Cls  cluster.Cluster
	Metr metrics.MetricsService
	Auth core.Authenticator
	Sql  persistence.Postgresql
}

func (s *AppManager) Metrics() metrics.MetricsService {
	return s.Metr
}
func (s *AppManager) Cluster() cluster.Cluster {
	return s.Cls
}
func (s *AppManager) Authenticator() core.Authenticator {
	return s.Auth
}

func (s *AppManager) Start(f conf.Env, c cluster.Cluster) error {
	s.Cls = c
	tkn := util.JwtHMac{Alg: "SHS256", Ksz: core.JWT_KEY_SIZE}
	ci := util.Aes{Ksz: core.CIPHER_KEY_SIZE}
	err := c.Atomic(func(ctx cluster.Ctx) error {
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
	s.Auth = &AuthManager{Tkn: &tkn, Cipher: &ci, Kid: f.GroupName}
	sql := persistence.Postgresql{Url: f.Pgs.DatabaseURL}
	err = sql.Create()
	if err != nil {
		return err
	}
	s.Sql = sql
	ms := persistence.MetricsDB{Sql: &sql}
	s.Metr = &ms
	return nil
}

func (s *AppManager) Shutdown() {
	s.Sql.Close()
}

func (s *AppManager) Create(classId int) event.Event {
	return &event.Login{}
}

func (s *AppManager) OnEvent(e event.Event) {

}
