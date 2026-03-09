package main

import (
	"gameclustering.com/internal/core"
)

type LocalEventListener struct {
	*PostofficeService
}

func (s *LocalEventListener) OnError(e core.Event, err error) {
	core.AppLog.Printf("On event error %v %s\n", e, err.Error())
}

func (s *LocalEventListener) OnEvent(e core.Event) {
	core.AppLog.Printf("On event %v\n", e)
}
