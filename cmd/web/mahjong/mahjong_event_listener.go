package main

import (
	"fmt"

	"gameclustering.com/internal/event"
)

type MahjongEventListener struct {
	*MahjongService
}

func (s *MahjongEventListener) OnError(e event.Event, err error) {
	fmt.Printf("On event error %v %s\n", e, err.Error())
}

func (s *MahjongEventListener) OnEvent(e event.Event) {
	fmt.Printf("On event %v\n", e)
}
