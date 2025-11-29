package main

import (
	"time"

	"gameclustering.com/internal/event"
)

type MahjongTimeout interface {
	Start(t *MahjongTable)
	Stop(commited bool)
	OId() int64
}

type MahjongEventObj struct {
	event.EventObj
	SystemId int64
	TableId  int64
}

func (s *MahjongEventObj) RecipientId() int64 {
	return s.SystemId
}

func (s *MahjongEventObj) OnRecipientId(recipientId int64) {
	s.SystemId = recipientId
}

type OnTurn func()

type MahjongTimeoutObj struct {
	MahjongEventObj
	Commited chan bool
	N        MahjongPlayTurn
	T        OnTurn //triger on timer
	P        OnTurn //call on stop
}

func (s *MahjongTimeoutObj) Start(tb *MahjongTable) {
	tm := *time.NewTimer(time.Duration(s.N.CountDown+COUNT_DOWN_BUFFER) * time.Second)
	s.Commited = make(chan bool)
	closing := false
	defer close(s.Commited)
	for {
		if closing {
			break
		}
		select {
		case <-tm.C:
			s.T()
		case c := <-s.Commited:
			if c && s.P != nil {
				s.P()
			}
			closing = true
		}
	}
}

func (mt *MahjongTimeoutObj) Stop(commited bool) {
	mt.Commited <- commited
}
