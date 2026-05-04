package main

import (
	"context"
	"fmt"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/postoffice/clustering"
)

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
	s.AppManager.Start(env)
	s.createSchema()

	m := clustering.MemberlistManager{StoreDir: fmt.Sprintf("%s/%s", env.HomeDir, env.GroupName)}
	m.Seed = []string{"192.168.1.11", "192.168.1.3"}
	m.Binding = env.NodeName
	err := m.Start(fmt.Appendf([]byte{}, "%s:%s", s.Context(), s.NodeId()), s.Sequence())
	if err != nil {
		core.AppLog.Printf("no cluster can join %s", err.Error())
		return err
	}
	s.mm = &m
	s.mm.DWait.Wait()
	s.started = true
	ak, err := m.AuthKey(context.Background(), &protocol.Request{Context: s.F.PresenceCtx()})
	if err != nil {
		panic(err.Error())
	}
	au, err := s.LoadAuth(ak)
	if err != nil {
		panic(err.Error())
	}
	s.Auth = au
	core.AppLog.Debug().Msgf("postoffice service started %s %s", env.HttpBinding, env.HomeDir)
	return nil
}

func (s *PostofficeService) Shutdown() {
	s.started = false
	core.AppLog.Debug().Msg("postoffice service shutting down ...")
	s.AppManager.Shutdown()
	s.mm.ShutdownHook()
}
