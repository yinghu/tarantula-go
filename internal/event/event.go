package event

import (
	"gameclustering.com/internal/core"
)

type Chunk struct {
	Remaining bool
	Data      []byte
}

type EventFactory interface {
	Create(classId int) Event
}

type Event interface {
	OnTopic() bool
	Streaming(c Chunk)
	core.Persistentable
}

type EventObj struct {
	Topic    bool
	Listener chan Chunk
}

func (s *EventObj) OnTopic() bool {
	return s.Topic
}

func (s *EventObj) Streaming(c Chunk) {
	s.Listener <- c
}
