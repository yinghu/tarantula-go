package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type LocalEventListener struct {
	*PostofficeService
}

func (s *LocalEventListener) OnError(e event.Event, err error) {
	core.AppLog.Printf("On event error %v %s\n", e, err.Error())
}

func (s *LocalEventListener) OnEvent(e event.Event) {
	//core.AppLog.Printf("On event %v\n", e)
}
