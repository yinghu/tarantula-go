package main

import (
	"gameclustering.com/internal/core"
)

type MahjongEventListener struct {
	*MahjongService
}

func (s *MahjongEventListener) OnError(e core.Event, err error) {
	core.AppLog.Printf("On event error %v %s\n", e, err.Error())
}

func (s *MahjongEventListener) OnEvent(e core.Event) {
	ex, y := e.(*MahjongEvent)
	if !y {
		return
	}
	ex.Token.TableId = ex.TableId
	s.Dispatcher <- ex.Token
}
