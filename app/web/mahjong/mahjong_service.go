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
		core.AppLog.Printf("joined from %d\n", e.RecipientId())
		s.Dispatcher <- MahjongPlayToken{SystemId: e.RecipientId(), Cmd: CMD_JOINED}
	case event.KICKOFF_CID:
		core.AppLog.Printf("kickoff from %d\n", e.RecipientId())
		s.Dispatcher <- MahjongPlayToken{SystemId: e.RecipientId(), Cmd: CMD_LEFT}
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
			s.onTable(t.SystemId)
		case CMD_LEFT:
			s.offTable(t.SystemId)
		}
	}
}

func (s *MahjongService) onTable(systemId int64) {
	tid, _ := s.Sequence().Id()
	table := MahjongTable{Id: tid, MahjongService: s, Setup: mj.ClassicMahjong{}}
	table.Reset()
	s.TableIndex[systemId] = &table
	go table.Play()
	mt := MahjongTableEvent{TableId: table.Id, SystemId: systemId}
	s.Pusher().Push(&mt)
}
func (s *MahjongService) offTable(systemId int64) {
	table, exists := s.TableIndex[systemId]
	if !exists {
		return
	}
	delete(s.TableIndex, systemId)
	table.Turn <- MahjongPlayToken{Cmd: CMD_END}

}
