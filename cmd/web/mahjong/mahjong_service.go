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
	ClassicMahjong
}

func (s *MahjongService) Config() string {
	return "/etc/tarantula/mahjong-conf.json"
}

func (s *MahjongService) Start(f conf.Env, c core.Cluster, p event.Pusher) error {
	s.ItemUpdater = s
	s.AppManager.Start(f, c, p)
	s.ClassicMahjong = ClassicMahjong{}
	s.ClassicMahjong.New()
	http.Handle("/mahjong/dice", bootstrap.Logging(&MahjongDicer{MahjongService: s}))
	http.Handle("/mahjong/claim", bootstrap.Logging(&MahjongClaimer{MahjongService: s}))
	return nil
}


func (s *MahjongService) Create(classId int, topic string) (event.Event, error) {
	me := event.CreateEvent(classId)
	me.OnListener(&MahjongEventListener{})//inbound event callback
	me.OnTopic(topic)
	if me == nil {
		return nil, fmt.Errorf("event ( %d ) not supported", classId)
	}
	return me, nil
}

func (s *MahjongService) VerifyTicket(ticket string) error {
	core.AppLog.Printf("validate ticket %s\n", ticket)
	//_, err := s.auth.ValidateTicket(ticket)
	//if err != nil {
	//return err
	//}
	return nil
}

func (s *MahjongService) OnError(e event.Event, err error) {
	core.AppLog.Printf("On event error %s\n", err.Error())
}

func (s *MahjongService) OnEvent(e event.Event) {
	switch e.ClassId() {
	case event.MESSAGE_CID:
		s.Pusher().Push(e)
	default:

	}
	//se, isSe := e.(*event.SubscriptionEvent)
	//if isSe {
	//for i := range s.topicQ {
	//s.topicQ[i] <- *se
	//}
	//return
	//}
	//s.inboundQ <- e
}
