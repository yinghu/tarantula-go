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
	Create(classId int, topic string) (Event, error)
	VerifyTicket(ticket string) error
	EventListener
	Postoffice
}

type Event interface {
	Inbound(buff core.DataBuffer) error
	Outbound(buff core.DataBuffer) error
	core.Persistentable
	Listener() EventListener
	Topic(t string)
	OnTopic() string
	OnIndex(ds core.DataStore) error
}

type Postoffice interface {
	Send(e Event) error
	List(q Query)
	Recover(q Query)
}

type EventObj struct {
	Callback EventListener `json:"-"`
	core.PersistentableObj
	topic string `json:"-"`
}

func (s *EventObj) Topic(t string) {
	s.topic = t
}

func (s *EventObj) OnTopic() string {
	return s.topic
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

func (s *EventObj) OnIndex(ds core.DataStore) error {
	return nil
}
