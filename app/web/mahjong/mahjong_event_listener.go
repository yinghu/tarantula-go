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
	ex, y := e.(*MahjongEvent)
	if !y {
		return
	}
	core.AppLog.Printf("On event %v\n", ex)
	dice := s.Table.Setup.Dice()
	me := MahjongDiceEvent{Dice1: int32(dice[0]), Dice2: int32(dice[1])}
	//s.Table.Setup.Draw()
	s.Pusher().Push(&me)
	mt :=  MahjongTableEvent{Table:&s.Table}
	s.Pusher().Push(&mt)

}
