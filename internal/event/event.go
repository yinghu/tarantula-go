package event

import (
	"gameclustering.com/internal/persistence"
)

type Chunk struct {
	Remaining bool
	Data      []byte
}

type Event interface {
	OnTopic() bool
	Streaming(c Chunk)
	persistence.Persistentable
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
