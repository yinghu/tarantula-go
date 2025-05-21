package main

import (
	"gameclustering.com/internal/event"
)

type PresenceFactory struct {
}

func (s *PresenceFactory) Create(classId int) event.Event {
	return &Login{}
}
