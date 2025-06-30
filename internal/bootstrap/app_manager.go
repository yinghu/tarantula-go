package bootstrap

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/metrics"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type AppManager struct {
	cls        core.Cluster
	metr       metrics.MetricsService
	imse       item.ItemService
	auth       core.Authenticator
	Sql        persistence.Postgresql
	ctx        string
	standalone bool
	appAuth    core.Authenticator
	seq        core.Sequence
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
func (s *AppManager) Sequence() core.Sequence {
	return s.seq
}

func (s *AppManager) Start(f conf.Env, c core.Cluster) error {
	core.AppLog.Printf("app manager starting %v\n", f)
	s.cls = c
	s.ctx = f.GroupName
	s.standalone = f.Standalone
	sfk := util.NewSnowflake(f.NodeId, util.EpochMillisecondsFromMidnight(2020, 1, 1))
	s.seq = &sfk
	au, err := s.LoadAuth(f.Presence, f.GroupName)
	if err != nil {
		return nil
	}
	s.auth = au
	ap, err := s.LoadAuth("presence", "presence")
	if err != nil {
		return err
	}
	s.appAuth = ap
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

	c.Started()
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

func (s *AppManager) MemberJoined(joined core.Node) {
	if s.standalone {
		return
	}
	core.AppLog.Printf("Member joined %v\n", joined)
	tick, err := s.appAuth.CreateTicket(1, 1, SUDO_ACCESS_CONTROL)
	if err != nil {
		return
	}
	data, err := json.Marshal(joined)
	if err != nil {
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://web/asset/admin/join", bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tick)
	resp, err := client.Do(req)
	if err != nil {
		core.AppLog.Printf("Error %s\n", err.Error())
		return
	}
	defer resp.Body.Close()
	r, err := io.ReadAll(resp.Body)
	if err != nil {
		core.AppLog.Printf("Resp Error %s\n", err.Error())
	}
	core.AppLog.Printf("Response code : %d %s\n", resp.StatusCode, string(r))
}

func (s *AppManager) Context() string {
	return s.ctx
}

func (s *AppManager) Service() TarantulaService {
	return s
}

func (s *AppManager) LoadAuth(context string, group string) (core.Authenticator, error) {
	tkn := util.JwtHMac{Alg: core.JWT_ALG, Ksz: core.JWT_KEY_SIZE}
	ci := util.Aes{Ksz: core.CIPHER_KEY_SIZE}
	err := s.cls.Atomic(context, func(ctx core.Ctx) error {
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
		return nil, err
	}
	err = s.cls.Atomic(context, func(ctx core.Ctx) error {
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
		return nil, err
	}
	return &AuthManager{Tkn: &tkn, Cipher: &ci, Kid: group}, nil
}
