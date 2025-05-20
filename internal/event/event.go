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
	Inbound(buff core.DataBuffer)
	Outbound(buff core.DataBuffer)
	core.Persistentable
}



type EventService interface {
	Publish(e Event) error
}

type EventObj struct {
	Topic    bool
	Listener chan Chunk
	core.PersistentableObj
}

func (s *EventObj) OnTopic() bool {
	return s.Topic
}

func (s *EventObj) Inbound(buff core.DataBuffer) {

}

func (s *EventObj) Outbound(buff core.DataBuffer) {

}
