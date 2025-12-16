package main

import (
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/mj"
)

type MahjongService struct {
	bootstrap.AppManager
	TableIndex map[int64]*MahjongTable
	Dispatcher chan MahjongPlayToken
}

func (s *MahjongService) Config() string {
	return "/etc/tarantula/mahjong-conf.json"
}

func (s *MahjongService) Start(f conf.Env, c core.Cluster, p event.Pusher) error {
	s.ItemUpdater = s
	s.AppManager.Start(f, c, p)
	s.TableIndex = make(map[int64]*MahjongTable)
	s.Dispatcher = make(chan MahjongPlayToken, 10)
	go s.dispatch()
	http.Handle("/mahjong/table", bootstrap.Logging(&MahjongTableSelector{MahjongService: s}))
	return nil
}

func (s *MahjongService) Shutdown() {
	core.AppLog.Println("majong service shutting down ...")
	close(s.Dispatcher)
	s.AppManager.Shutdown()
}

func (s *MahjongService) Create(classId int, topic string) (event.Event, error) {
	e := event.CreateEvent(classId)
	if e != nil {
		e.OnTopic(topic)
		e.OnListener(s)
		return e, nil
	}
	me := MahjongEvent{}
	me.OnListener(&MahjongEventListener{MahjongService: s})
	return &me, nil
}
func (s *MahjongService) VerifyTicket(ticket string) (core.OnSession, error) {
	session, err := s.AppAuth.ValidateTicket(ticket)
	if err != nil {
		return session, err
	}
	if session.AccessControl < bootstrap.PROTECTED_ACCESS_CONTROL {
		return session, fmt.Errorf("player access control required %d", session.AccessControl)
	}
	return session, nil
}

func (s *MahjongService) OnError(e event.Event, err error) {
	core.AppLog.Printf("On event error %s\n", err.Error())
}

func (s *MahjongService) OnEvent(e event.Event) {
	switch e.ClassId() {
	case event.MESSAGE_CID:
		s.Pusher().Push(e)
	case event.JOIN_CID:
		join, _ := e.(*event.JoinEvent)
		core.AppLog.Printf("joined from %d %d\n", join.RecipientId(), join.Flag)
		s.Dispatcher <- MahjongPlayToken{SystemId: e.RecipientId(), Cmd: CMD_JOINED, TableId: join.Flag}
	case event.KICKOFF_CID:
		kickoff, _ := e.(*event.KickoffEvent)
		core.AppLog.Printf("kickoff from %d : %d\n", e.RecipientId(), kickoff.Flag)
		s.Dispatcher <- MahjongPlayToken{SystemId: e.RecipientId(), Cmd: CMD_LEFT, TableId: kickoff.Flag}
		id, _ := s.Sequence().Id()
		e.OnOId(id)
		e.OnTopic("mahjong")
		err := s.Send(e)
		if err != nil {
			core.AppLog.Printf("failed to send event %s\n", err.Error())
		}
	default:
	}
}

func (s *MahjongService) dispatch() {
	for t := range s.Dispatcher {
		switch t.Cmd {
		case CMD_JOINED:
			s.onTable(t.SystemId, t.TableId)
		case CMD_LEFT:
			s.offTable(t.SystemId, t.TableId)
		}
	}
}

func (s *MahjongService) onTable(systemId int64, flag int64) {
	core.AppLog.Printf("table flag %d\n", flag)
	if systemId != flag {
		t, exists := s.TableIndex[flag]
		if !exists {
			//tid, _ := s.Sequence().Id()
			t = &MahjongTable{Id: flag, Pusher: s.Pusher(), Sequence: s.Sequence(), CMJ: mj.ClassicMahjong{}, Solo: flag == systemId}
			s.TableIndex[flag] = t
			go t.Play()
		}
		pt := MahjongPlayToken{SystemId: systemId, Cmd: CMD_SIT}
		t.Sync <- pt
		core.AppLog.Printf("joining table flag %d\n", flag)
		return
	}
	tid, _ := s.Sequence().Id()
	table := MahjongTable{Id: tid, Pusher: s.Pusher(), Sequence: s.Sequence(), CMJ: mj.ClassicMahjong{}, Solo: flag == systemId}
	table.New()
	s.TableIndex[flag] = &table
	go table.Play()
	pt := MahjongPlayToken{SystemId: systemId, Cmd: CMD_SIT}
	table.Sync <- pt
}
func (s *MahjongService) offTable(systemId int64, tableId int64) {
	core.AppLog.Printf("table id : %d\n", tableId)
	table, exists := s.TableIndex[systemId]
	if !exists {
		return
	}
	delete(s.TableIndex, systemId)
	table.Sync <- MahjongPlayToken{Cmd: CMD_END}

}
