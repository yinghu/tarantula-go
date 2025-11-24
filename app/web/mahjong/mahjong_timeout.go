package main

import (
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

const (
	COUNT_DOWN_BUFFER int = 5
)

type MahjongTimeout interface {
	Start(t *MahjongTable)
	Stop()
	OId() int64
}

type OnTurn func()

type MahjongTimeoutObj struct {
	event.EventObj
	Commited chan bool
	N        MahjongPlayTurn
	T        OnTurn
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
		case <-s.Commited:
			core.AppLog.Printf("Stop %d\n", s.OId())
			closing = true
		}
	}
}

func (mt *MahjongTimeoutObj) Stop() {
	mt.Commited <- true
}
