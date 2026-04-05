package main

import (
	"context"
	"fmt"
	"time"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/postoffice/clustering"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LogData struct {
	level zerolog.Level
	log   []byte
}

type PostofficeService struct {
	bootstrap.AppManager
	mm      *clustering.MemberlistManager
	mLog    chan LogData
	started bool
}

func (s *PostofficeService) Config() string {
	return "/etc/tarantula/postoffice-conf.json"
}

func (s *PostofficeService) Start(env core.Env) error {
	env.AuthLevel = core.ADMIN_ACCESS_CONTROL
	env.IsClusterMember = true
	s.mLog = make(chan LogData, 100)
	s.RegisterLogForwarder(s)
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
	go s.runForward()
	core.AppLog.Debug().Msgf("postoffice service started %s %s", env.HttpBinding, env.HomeDir)
	return nil
}

func (s *PostofficeService) Shutdown() {
	s.started = false
	core.AppLog.Debug().Msg("postoffice service shutting down ...")
	s.AppManager.Shutdown()
	s.mm.ShutdownHook()
}

func (s *PostofficeService) Forward(level zerolog.Level, log []byte) {
	s.mLog <- LogData{level: level, log: log}
}

func (s *PostofficeService) runForward() {
	time.Sleep(3 * time.Second)
	for s.started {
		for data := range s.mLog {
			lf := event.LogEventFactory{}
			e := protocol.LogEvent{}
			err := protojson.Unmarshal(data.log, &e)
			if err != nil {
				e.Level = "error"
				e.Message = err.Error()
				e.Time = timestamppb.Now()
				e.Source = "postoffice:64"
			}
			id, _ := s.Sequence().Id()
			t, _ := lf.FromLogEvent(&e)
			t.NodeId = s.NodeId()
			t.Tag = s.Context()
			t.Event.Id = uint64(id)
			s.mm.DataServiceProvider.Publish(context.Background(), t)
		}
	}
	close(s.mLog)
}
