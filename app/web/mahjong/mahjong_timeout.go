package main

import (
	"sync"
	"time"

	"gameclustering.com/internal/core"
)

type MahjongTimeout interface {
	Start(t *MahjongTable)
	Stop(commited MahjongPlayToken, closing bool)
	OId() int64
}

type MahjongEventObj struct {
	core.EventObj
	SystemId int64
	TableId  int64
}

func (s *MahjongEventObj) RecipientId() int64 {
	return s.SystemId
}

func (s *MahjongEventObj) OnRecipientId(recipientId int64) {
	s.SystemId = recipientId
}

type OnTimer func()
type OnStop func(commited MahjongPlayToken)

type MahjongTimeoutObj struct {
	MahjongEventObj
	N       MahjongPlayTurn
	T       OnTimer //triger on timer
	P       OnStop  //call on stop
	K       *time.Timer
	Lock    *sync.Mutex
	Stopped bool
}
type StopSignal struct {
	Commited bool
	Closing  bool
}

func (s *MahjongTimeoutObj) Start(tb *MahjongTable) {
	s.Lock = &sync.Mutex{}
	s.Stopped = false
	s.K = time.AfterFunc(time.Duration(s.N.CountDown+COUNT_DOWN_BUFFER)*time.Second, func() {
		s.Lock.Lock()
		defer s.Lock.Unlock()
		if s.Stopped {
			return
		}
		s.T()
	})
}

func (s *MahjongTimeoutObj) Stop(commited MahjongPlayToken, closing bool) {
	s.Lock.Lock()
	s.K.Stop()
	defer s.Lock.Unlock()
	s.Stopped = true
	if closing {
		core.AppLog.Printf("timeout closed forcefully %d", s.OId())
		return
	}
	s.P(commited)
}
