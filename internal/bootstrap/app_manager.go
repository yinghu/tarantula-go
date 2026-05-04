package bootstrap

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
	"github.com/rs/zerolog"
)

type AppManager struct {
	imse        core.ItemService
	Auth        core.Authenticator
	Sql         persistence.Postgresql
	F           core.Env
	seq         core.Sequence
	ItemUpdater core.ItemListener
	tcpPusher   core.Pusher
	ManagedApps []string
	cluster     *ClusterManager
	event       *EventManager
	log         io.Writer //zerolog.Logger
	forward     LogForwarder
	threshold   zerolog.Level
}

func (s *AppManager) ItemService() core.ItemService {
	return s.imse
}

func (s *AppManager) Pusher() core.Pusher {
	return s.tcpPusher
}

func (s *AppManager) Authenticator() core.Authenticator {
	return s.Auth
}
func (s *AppManager) Sequence() core.Sequence {
	return s.seq
}
func (s *AppManager) ItemListener() core.ItemListener {
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

func (s *AppManager) RegisterLogForwarder(threshold zerolog.Level, logf LogForwarder) {
	s.forward = logf
}
func (s *AppManager) Start(f core.Env) error {
	s.F = f
	s.initLogger(f)
	core.AppLog.Info().Msgf("app manager starting on %s %v\n", f.Prefix, f)
	s.event = &EventManager{App: s}
	s.ManagedApps = f.ManagedApps
	sfk := util.NewSnowflake(f.NodeId, util.EpochMillisecondsFromMidnight(2020, 1, 1))
	s.seq = &sfk
	if f.Pgs.Enabled {
		core.AppLog.Info().Msgf("connecting sql %s",f.Pgs.DatabaseURL)
		dbCreate := persistence.Postgresql{Url: f.Pgs.DatabaseURL + "/postgres"}
		err := dbCreate.CreateDatabase(fmt.Sprintf("CREATE DATABASE %s_%s_%s", f.Prefix, "tarantula", f.GroupName))
		if err != nil {
			core.AppLog.Warn().Msgf("failed to create database %s", err.Error())
		}
		sql := persistence.Postgresql{Url: f.Pgs.DatabaseURL + "/" + f.Prefix + "_tarantula_" + f.GroupName}
		err = sql.Create()
		if err != nil {
			return err
		}
		s.Sql = sql
		gitStore := persistence.GitItemStore{RepositoryDir: f.HomeDir + "/bin/tarantula", JsonRequester: s}
		gitStore.Start()
		is := persistence.ItemDB{Sql: &sql, Gis: &gitStore}
		err = is.Start()
		if err != nil {
			return err
		}
		s.imse = &is
	}
	if f.IsClusterMember {
		return nil
	}
	core.AppLog.Warn().Msgf("Starting cluster client to %s", f.Host)
	s.cluster = &ClusterManager{App: s}
	s.RegisterLogForwarder(zerolog.DebugLevel, s.cluster)

	err := s.cluster.connect(f.Host)
	if err != nil {
		panic(err.Error())
	}
	ak, err := s.cluster.AuthKey()
	if err != nil {
		panic(err.Error())
	}
	au, err := s.LoadAuth(ak)
	if err != nil {
		panic(err.Error())
	}
	s.Auth = au
	return nil
}

func (s *AppManager) Shutdown() {
	if !s.F.IsClusterMember {
		s.cluster.disconnect()
	}
	util.GitPush()
	s.Sql.Close()
	core.AppLog.Info().Msg("app manager shutting down ...")
}

func (s *AppManager) Context() string {
	return s.F.GroupName
}

func (s *AppManager) Service() TarantulaService {
	return s
}

func (s *AppManager) LoadAuth(ak *protocol.AuthKey) (core.Authenticator, error) {
	tkn := util.JwtHMac{Alg: core.JWT_ALG, Ksz: core.JWT_KEY_SIZE}
	ci := util.Aes{Ksz: core.CIPHER_KEY_SIZE}
	jk, err := util.KeyFromBase64(string(ak.Jwt))
	if err != nil {
		return nil, err
	}
	tkn.HMacFromKey(jk)
	ck, err := util.KeyFromBase64(string(ak.Cipher))
	if err != nil {
		return nil, err
	}
	ci.AesGcmFromKey(ck)
	return &AuthManager{Tkn: &tkn, Cipher: &ci, Kid: "presence"}, nil
}

func (c *AppManager) Write(data []byte) (int, error) {
	return c.log.Write(data)
}

func (c *AppManager) WriteLevel(level zerolog.Level, data []byte) (int, error) {
	if c.forward != nil && level >= c.threshold {
		cp := append([]byte{}, data...)
		c.forward.Forward(level, cp)
	}
	return c.Write(data)
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
