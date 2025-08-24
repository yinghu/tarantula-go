package bootstrap

import (
	"fmt"
	"time"

	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/metrics"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type AppManager struct {
	cls         core.Cluster
	metr        metrics.MetricsService
	imse        item.ItemService
	auth        core.Authenticator
	Bsl         BootstrapListener
	Sql         persistence.Postgresql
	ctx         string
	standalone  bool
	AppAuth     core.Authenticator
	seq         core.Sequence
	ItemUpdater item.ItemListener
	ManagedApps []string
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
func (s *AppManager) ItemListener() item.ItemListener {
	return s.ItemUpdater
}
func (s *AppManager) BootstrapListener() BootstrapListener {
	return s.Bsl
}
func (s *AppManager) Start(f conf.Env, c core.Cluster) error {
	core.AppLog.Printf("app manager starting on %s %v\n", f.Prefix, f)
	s.ManagedApps = f.ManagedApps
	s.cls = c
	s.ctx = f.GroupName
	s.standalone = f.Standalone
	sfk := util.NewSnowflake(f.NodeId, util.EpochMillisecondsFromMidnight(2020, 1, 1))
	s.seq = &sfk
	fctx := f.PresenceCtx()
	if f.GroupName == "admin" {
		fctx = f.ClusterCtx()
	}
	au, err := s.LoadAuth(fctx)
	if err != nil {
		return nil
	}
	s.auth = au
	ap, err := s.LoadAuth(f.PresenceCtx())
	if err != nil {
		return err
	}
	s.AppAuth = ap
	sql := persistence.Postgresql{Url: f.Pgs.DatabaseURL + "/" + f.Prefix + "_tarantula_" + f.GroupName}
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
	gitStore := persistence.GitItemStore{RepositoryDir: f.HomeDir + "/bin/tarantula", JsonRequester: s}
	gitStore.Start()
	is := persistence.ItemDB{Sql: &sql, Gis: &gitStore, Cls: s.cls}
	err = is.Start()
	if err != nil {
		return err
	}
	s.imse = &is
	c.Started()
	return nil
}

func (s *AppManager) Shutdown() {
	util.GitPush()
	s.Sql.Close()
}

func (s *AppManager) Create(classId int, topic string) (event.Event, error) {
	return nil, nil
}

func (s *AppManager) VerifyTicket(ticket string) error {
	_, err := s.auth.ValidateTicket(ticket)
	if err != nil {
		return err
	}
	return nil
}
func (s *AppManager) Send(e event.Event) error {
	for i := range 5 {
		ret := s.PostJsonSync(fmt.Sprintf("%s/%s/%d", "http://postoffice:8080/postoffice/publish", e.Topic(), e.ClassId()), e)
		if ret.ErrorCode == 0 {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
		core.AppLog.Printf("Retries: %d %v\n", i, ret)
	}
	return fmt.Errorf("failed after retries")
}
func (s *AppManager) List(query event.Query) {
	s.PostJsonAsync(fmt.Sprintf("%s/%d", "http://postoffice:8080/postoffice/query", query.QId()), query, query.QCc())
}
func (s *AppManager) Recover(query event.Query) {
	for i := range 5 {
		ret := s.PostJsonSync(fmt.Sprintf("%s/%d", "http://postoffice:8080/postoffice/recover", query.QId()), query)
		if ret.ErrorCode == 0 {
			return
		}
		time.Sleep(1000 * time.Millisecond)
		core.AppLog.Printf("Retries: %d %v\n", i, ret)
	}
}

func (s *AppManager) OnEvent(e event.Event) {

}

func (s *AppManager) OnError(e event.Event, err error) {

}

func (s *AppManager) Context() string {
	return s.ctx
}

func (s *AppManager) Service() TarantulaService {
	return s
}

func (s *AppManager) NodeStarted(n core.Node) {
	core.AppLog.Printf("Node started %s\n", n.Name)
}

func (s *AppManager) LoadAuth(context string) (core.Authenticator, error) {
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
	return &AuthManager{Tkn: &tkn, Cipher: &ci, Kid: "presence"}, nil
}
