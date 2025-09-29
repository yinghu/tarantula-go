package main

import (

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type MahjongEventListener struct {
	*MahjongService
}

func (s *MahjongEventListener) OnError(e event.Event, err error) {
	core.AppLog.Printf("On event error %v %s\n", e, err.Error())
}

func (s *MahjongEventListener) OnEvent(e event.Event) {
	core.AppLog.Printf("On event %v\n", e)

}
