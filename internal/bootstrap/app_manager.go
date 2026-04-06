package bootstrap

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
	"github.com/rs/zerolog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"google.golang.org/grpc"
)

type AppManager struct {
	metr        core.MetricsService
	imse        item.ItemService
	auth        core.Authenticator
	Sql         persistence.Postgresql
	F           core.Env
	seq         core.Sequence
	ItemUpdater item.ItemListener
	tcpPusher   core.Pusher
	ManagedApps []string
	cluster     *ClusterManager
	event       *EventManager
	log         io.Writer //zerolog.Logger
	forward     LogForwarder
}

func (s *AppManager) ItemService() item.ItemService {
	return s.imse
}

func (s *AppManager) Pusher() core.Pusher {
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

func (c *AppManager) Cluster() core.ClusterService {
	return c.cluster
}
func (c *AppManager) Event() core.EventService {
	return c.event
}

func (s *AppManager) NodeId() string {
	return s.F.NodeName
}

func (s *AppManager) ClusterMember() bool {
	return s.F.IsClusterMember
}

func (s *AppManager) RegisterLogForwarder(logf LogForwarder) {
	s.forward = logf
}
func (s *AppManager) Start(f core.Env) error {
	s.F = f
	s.initLogger(f)
	//CreateAppLog(f.LogDir, f.LogTruncated, f.Standalone, s)
	core.AppLog.Printf("app manager starting on %s %v\n", f.Prefix, f)
	s.event = &EventManager{App: s}
	s.ManagedApps = f.ManagedApps
	sfk := util.NewSnowflake(f.NodeId, util.EpochMillisecondsFromMidnight(2020, 1, 1))
	s.seq = &sfk
	fctx := f.PresenceCtx()
	au, err := s.LoadAuth(fctx)
	if err != nil {
		return nil
	}
	s.auth = au
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
	if f.IsClusterMember {
		return nil
	}
	core.AppLog.Warn().Msgf("Starting cluster client to %s", f.Host)
	s.cluster = &ClusterManager{App: s}
	s.RegisterLogForwarder(s.cluster)
	task := make(chan core.Task, 100)
	s.cluster.wTask = task
	tex := core.ServiceCallOperator{RTask: task, LocalConns: make(map[string]*grpc.ClientConn)}
	go tex.RunTask()
	return s.cluster.connect(f.Host)
}

func (s *AppManager) Shutdown() {
	if !s.F.IsClusterMember {
		s.cluster.disconnect()
	}
	util.GitPush()
	s.Sql.Close()
	core.AppLog.Println("app manager shutting down ...")
}

func (s *AppManager) Context() string {
	return s.F.GroupName
}

func (s *AppManager) Service() TarantulaService {
	return s
}

func (s *AppManager) Forward(topic *protocol.Topic) {
	fmt.Printf("topic forward process %v\n", topic)
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

func (c *AppManager) Write(data []byte) (int, error) {
	fmt.Printf("LOG : %s\n", string(data))
	return len(data), nil
}

func (c *AppManager) WriteLevel(level zerolog.Level, data []byte) (int, error) {
	if c.forward != nil {
		cp := append([]byte{}, data...)
		c.forward.Forward(level, cp)
	}
	return c.log.Write(data)
}

func (c *AppManager) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	_, f, line, ok := runtime.Caller(3)
	if !ok {
		e.Str("source", "unknown")
		return
	}
	e.Str("source", fmt.Sprintf("%s:%d", f, line))
}

func (c *AppManager) initLogger(f core.Env) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = time.RFC3339
	if f.Standalone {
		CreateTestLog()
		return
	}
	err := os.MkdirAll(f.LogDir+"/log", 0755)
	if err != nil {
		CreateTestLog()
		return
	}
	opt := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	if f.LogTruncated {
		opt = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	file, err := os.OpenFile(f.LogDir+"/log/tarantula.log", opt, 0644)
	if err != nil {
		CreateTestLog()
		return
	}
	c.log = file
	core.AppLog = zerolog.New(zerolog.MultiLevelWriter(c)).With().Timestamp().Logger().Hook(c)
	core.AppLog.Info().Msg("Initialized app log")

}
