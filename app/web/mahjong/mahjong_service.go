package main

import (
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type MahjongService struct {
	bootstrap.AppManager
	Table MahjongTable
}

func (s *MahjongService) Config() string {
	return "/etc/tarantula/mahjong-conf.json"
}

func (s *MahjongService) Start(f conf.Env, c core.Cluster, p event.Pusher) error {
	s.ItemUpdater = s
	s.AppManager.Start(f, c, p)
	tid, _ := s.Sequence().Id()
	s.Table = MahjongTable{MahjongService: s, Turn: make(chan MahjongPlayToken, 10), Id: tid}
	s.Table.Reset()
	s.Table.Dice()
	s.Table.Deal()
	go s.Table.Play()
	http.Handle("/mahjong/table", bootstrap.Logging(&MahjongTableSelector{MahjongService: s}))
	return nil
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
	case event.KICKOFF_CID:
		core.AppLog.Printf("kickoff from %d\n", e.RecipientId())
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
