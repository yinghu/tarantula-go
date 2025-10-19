package event

import (
	"gameclustering.com/internal/core"
)

type IndexListener interface {
	LocalStore() core.DataStore
	Index(e Event)
}

type EventListener interface {
	OnEvent(e Event)
	OnError(e Event, err error)
}

type EventService interface {
	EventCreator
	VerifyTicket(ticket string) (core.OnSession, error)
	EventListener
	Postoffice
}

type Event interface {
	Inbound(buff core.DataBuffer) error
	Outbound(buff core.DataBuffer) error
	core.Persistentable
	OnListener(el EventListener)
	Listener() EventListener
	OnIndex(ix IndexListener)

	OnTopic(t string)
	Topic() string
	OId() int64
	OnOId(id int64)
	RecipientId() int64
}

type EventCreator interface {
	Create(classId int, topic string) (Event, error)
}

type Postoffice interface {
	Send(e Event) error
	List(q Query)
	Recover(q Query)
	Load(e Query)
}

type Publisher interface {
	Publish(e Event, ticket string) error
	Connect() error
	Close() error
}

type Pusher interface {
	Push(e Event)
}

type EventObj struct {
	Callback EventListener `json:"-"`
	core.PersistentableObj
	ETopic string `json:"ETopic"`
	EOid   int64  `json:"EOid"`
}

func (s *EventObj) OnTopic(t string) {
	s.ETopic = t
}

func (s *EventObj) Topic() string {
	return s.ETopic
}

func (s *EventObj) Inbound(buff core.DataBuffer) error {
	return nil
}
func (s *EventObj) Outbound(buff core.DataBuffer) error {
	return nil
}
func (s *EventObj) OnListener(el EventListener) {
	s.Callback = el
}
func (s *EventObj) Listener() EventListener {
	return s.Callback
}
func (s *EventObj) OnIndex(idx IndexListener) {

}
func (s *EventObj) OnOId(oid int64) {
	s.EOid = oid
}

func (s *EventObj) OId() int64 {
	return s.EOid
}

func (s *EventObj) RecipientId() int64 {
	return 0
}
