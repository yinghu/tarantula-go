package event

import (
	"gameclustering.com/internal/persistence"
)

type Chunk struct {
	Remaining bool
	Data      []byte
}

type Event interface {
	Topic() bool
	Streaming(c Chunk) error
	persistence.Persistentable
}

type EventObj struct {
	topic    bool
	Listener chan Chunk
}

func (s *EventObj) Topic() bool {
	return s.topic
}

func (s *EventObj) Streaming(c Chunk) error {
	s.Listener <- c
	return nil
}
