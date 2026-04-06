package main

import (
	"fmt"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/postoffice/clustering"
	"github.com/rs/zerolog"
)

type LogData struct {
	level zerolog.Level
	log   []byte
}

type PostofficeService struct {
	bootstrap.AppManager
	mm      *clustering.MemberlistManager
	started bool
}

func (s *PostofficeService) Config() string {
	return "/etc/tarantula/postoffice-conf.json"
}

func (s *PostofficeService) Start(env core.Env) error {
	env.AuthLevel = core.ADMIN_ACCESS_CONTROL
	env.IsClusterMember = true
	fwd := LogForwarder{App: s}
	s.RegisterLogForwarder(&fwd)
	s.AppManager.Start(env)

	s.createSchema()

	m := clustering.MemberlistManager{StoreDir: fmt.Sprintf("%s/%s", env.HomeDir, env.GroupName), Seq: s.Sequence()}
	m.Seed = []string{"192.168.1.11", "192.168.1.6"}
	m.Binding = env.NodeName
	err := m.Start()
	if err != nil {
		core.AppLog.Printf("no cluster can join %s", err.Error())
		return err
	}
	s.mm = &m
	s.started = true
	s.mm.DWait.Wait()
	core.AppLog.Debug().Msgf("postoffice service started %s %s", env.HttpBinding, env.HomeDir)
	return nil
}

func (s *PostofficeService) Shutdown() {
	s.started = false
	core.AppLog.Debug().Msg("postoffice service shutting down ...")
	s.AppManager.Shutdown()
	s.mm.ShutdownHook()
}

func (s *PostofficeService) Forward(topic *protocol.Topic) {
	fmt.Printf("postoffice topic forward process %v\n", topic)
}
