package main

import (
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type MahjongTimeout interface {
	Start(t *MahjongTable)
	Stop(commited bool, closing bool)
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
type OnStop func(commited bool)

type MahjongTimeoutObj struct {
	MahjongEventObj
	S chan StopSignal
	N MahjongPlayTurn
	T OnTurn //triger on timer
	P OnStop //call on stop
}
type StopSignal struct {
	Commited bool
	Closing  bool
}

func (s *MahjongTimeoutObj) Start(tb *MahjongTable) {
	tm := *time.NewTimer(time.Duration(s.N.CountDown+COUNT_DOWN_BUFFER) * time.Second)
	s.S = make(chan StopSignal)
	defer func() {
		tm.Stop()
		close(s.S)
		core.AppLog.Printf("timeout!!")
	}()
	for {
		select {
		case <-tm.C:
			s.T()
		case <-s.S:
			return
		}
	}
}

func (mt *MahjongTimeoutObj) Stop(commited bool, closing bool) {
	mt.S <- StopSignal{Commited: commited, Closing: closing}
	mt.P(commited)
}
