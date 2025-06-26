package event

import (
	"gameclustering.com/internal/core"
)

type Chunk struct {
	Remaining bool
	Data      []byte
}

type EventListener interface {
	OnEvent(e Event)
}

type EventService interface {
	Create(classId int, ticket string) (Event, error)
	EventListener
}

type Event interface {
	Inbound(buff core.DataBuffer)
	Outbound(buff core.DataBuffer)
	OnError(err error)
	core.Persistentable
	Listener() EventListener
}

type EventPublisher interface {
	Publish(e Event)
}

type EventObj struct {
	Cc chan Chunk
	Cb EventListener
	core.PersistentableObj
}

func (s *EventObj) Inbound(buff core.DataBuffer) {

}

func (s *EventObj) Outbound(buff core.DataBuffer) {

}

func (s *EventObj) OnError(err error) {

}

func (s *EventObj) Listener() EventListener {
	return s.Cb
}
