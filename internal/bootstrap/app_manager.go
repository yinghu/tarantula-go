package bootstrap

import (
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/metrics"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type AppManager struct {
	cls  core.Cluster
	metr metrics.MetricsService
	imse item.ItemService
	auth core.Authenticator
	Sql  persistence.Postgresql
}

func (s *AppManager) ItemService() item.ItemService {
	return s.imse
}

func (s *AppManager) Metrics() metrics.MetricsService {
	return s.metr
}
func (s *AppManager) Cluster() core.Cluster {
	return s.cls
}
func (s *AppManager) Authenticator() core.Authenticator {
	return s.auth
}

func (s *AppManager) Start(f conf.Env, c core.Cluster) error {
	s.cls = c
	tkn := util.JwtHMac{Alg: core.JWT_ALG, Ksz: core.JWT_KEY_SIZE}
	ci := util.Aes{Ksz: core.CIPHER_KEY_SIZE}
	err := c.Atomic(f.Presence, func(ctx core.Ctx) error {
		jsk, err := ctx.Get(core.JWT_KEY_NAME)
		if err != nil {
			core.AppLog.Println("Create new jwt key")
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
	err = c.Atomic(f.Presence, func(ctx core.Ctx) error {
		csk, err := ctx.Get(core.CIPHER_KEY_NAME)
		if err != nil {
			core.AppLog.Println("Create new cipher key")
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
	s.auth = &AuthManager{Tkn: &tkn, Cipher: &ci, Kid: f.GroupName}
	sql := persistence.Postgresql{Url: f.Pgs.DatabaseURL}
	err = sql.Create()
	if err != nil {
		return err
	}
	s.Sql = sql
	ms := persistence.MetricsDB{Sql: &sql}
	err = ms.Start()
	if err != nil {
		return err
	}
	s.metr = &ms
	is := persistence.ItemDB{Sql: &sql}
	err = is.Start()
	if err != nil {
		return err
	}
	s.imse = &is
	return nil
}

func (s *AppManager) Shutdown() {
	s.Sql.Close()
}

func (s *AppManager) Create(classId int, magicHeader string) (event.Event, error) {
	return &event.Login{}, nil
}

func (s *AppManager) OnEvent(e event.Event) {

}

func (s *AppManager) OnError(e error) {

}

func (s *AppManager) Updated(key string, value string) {
	core.AppLog.Printf("Key updated %s %s\n", key, value)
}
