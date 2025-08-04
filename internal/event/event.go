package event

import (
	"gameclustering.com/internal/core"
)

type Chunk struct {
	Remaining bool
	Data      []byte
}

type Topic struct {
	Id   int32  `json:"Id"`
	Name string `json:"Name"`
	App  string `json:"App"`
}

type EventListener interface {
	OnEvent(e Event)
	OnError(err error)
}

type EventService interface {
	Create(classId int, ticket string) (Event, error)
	EventListener
}

type Event interface {
	Inbound(buff core.DataBuffer) error
	Outbound(buff core.DataBuffer) error
	core.Persistentable
	Listener() EventListener
}

type EventPublisher interface {
	Publish(e Event)
}

type EventObj struct {
	Callback EventListener
	core.PersistentableObj
}

func (s *EventObj) Inbound(buff core.DataBuffer) error {
	return nil
}
func (s *EventObj) Outbound(buff core.DataBuffer) error {
	return nil
}
func (s *EventObj) Listener() EventListener {
	return s.Callback
}
