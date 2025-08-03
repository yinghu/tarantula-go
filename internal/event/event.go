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
	OnError(err error)
}

type EventService interface {
	Create(classId int, ticket string) (Event, error)
	EventListener
}

type Event interface {
	Inbound(buff core.DataBuffer)
	Outbound(buff core.DataBuffer)
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

func (s *EventObj) Outbound(buff core.DataBuffer) {
	err := s.WriteKey(buff)
	if err!=nil {
		s.Cb.OnError(err)
		return
	}
	err = s.Write(buff)
	if err!=nil {
		s.Cb.OnError(err)
		return
	}
	s.Cb.OnEvent(s)
}

func (s *EventObj) Inbound(buff core.DataBuffer) {
	err := s.ReadKey(buff)
	if err!=nil {
		s.Cb.OnError(err)
		return
	}
	err = s.Read(buff)
	if err !=nil{
		s.Cb.OnError(err)
		return
	}
}

func (s *EventObj) OnError(err error) {

}

func (s *EventObj) Listener() EventListener {
	return s.Cb
}
