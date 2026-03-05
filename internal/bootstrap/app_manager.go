package bootstrap

import (
	"context"
	"fmt"
	"time"

	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type AppManager struct {
	metr        core.MetricsService
	imse        item.ItemService
	auth        core.Authenticator
	Sql         persistence.Postgresql
	F           conf.Env
	AppAuth     core.Authenticator
	seq         core.Sequence
	ItemUpdater item.ItemListener
	tcpPusher   event.Pusher
	ManagedApps []string
}

func (s *AppManager) ItemService() item.ItemService {
	return s.imse
}

func (s *AppManager) Pusher() event.Pusher {
	return s.tcpPusher
}

func (s *AppManager) Metrics() core.MetricsService {
	return s.metr
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

func (s *AppManager) Name() string {
	return s.F.GroupName
}
func (s *AppManager) Start(f conf.Env, p event.Pusher) error {
	core.AppLog.Printf("app manager starting on %s %v\n", f.Prefix, f)
	s.ManagedApps = f.ManagedApps
	s.tcpPusher = p
	s.F = f
	sfk := util.NewSnowflake(f.NodeId, util.EpochMillisecondsFromMidnight(2020, 1, 1))
	s.seq = &sfk
	fctx := f.AuthCtx()
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
	dbCreate := persistence.Postgresql{Url: f.Pgs.DatabaseURL + "/postgres"}
	err = dbCreate.CreateDatabase(fmt.Sprintf("CREATE DATABASE %s_%s_%s", f.Prefix, "tarantula", f.GroupName))
	if err != nil {
		core.AppLog.Warn().Msgf("failed to create database %s", err.Error())
	}
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
	is := persistence.ItemDB{Sql: &sql, Gis: &gitStore}
	err = is.Start()
	if err != nil {
		return err
	}
	s.imse = &is
	return nil
}

func (s *AppManager) Shutdown() {
	util.GitPush()
	s.Sql.Close()
	core.AppLog.Println("app manager shutting down ...")
}

func (s *AppManager) Create(classId int, topic string) (event.Event, error) {
	return nil, nil
}

func (s *AppManager) VerifyTicket(ticket string) (core.OnSession, error) {
	session, err := s.auth.ValidateTicket(ticket)
	if err != nil {
		return session, err
	}
	if session.AccessControl < core.ADMIN_ACCESS_CONTROL {
		return session, fmt.Errorf("admin access control required %d", session.AccessControl)
	}
	return session, nil
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

func (s AppManager) Load(query event.Query) {
	e := query.QEvent()
	s.PostJsonAsync(fmt.Sprintf("%s/%d", "http://postoffice:8080/postoffice/load", e.ClassId()), e, query.QCc())
}

func (s *AppManager) OnEvent(e event.Event) {

}

func (s *AppManager) OnError(e event.Event, err error) {

}

func (s *AppManager) Context() string {
	return s.F.GroupName
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
	err := s.Atomic(context, func(ctx core.Ctx) error {
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
	err = s.Atomic(context, func(ctx core.Ctx) error {
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

func (c *AppManager) Atomic(prefix string, t core.Exec) error {
	if prefix == "" {
		prefix = c.F.GroupName
		core.AppLog.Printf("Reset Lock prefix %s\n", prefix)
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.F.EtcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	defer cli.Close()
	session, err := concurrency.NewSession(cli)
	if err != nil {
		return err
	}
	defer session.Close()
	mutex := concurrency.NewMutex(session, prefix+"#lock")
	ctx := context.Background()
	mutex.Lock(ctx)
	defer mutex.Unlock(ctx)
	return t(&core.EtcdClient{Cli: cli, Prefix: prefix})
}
