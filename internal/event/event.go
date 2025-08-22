package event

import (
	"gameclustering.com/internal/core"
)

type Chunk struct {
	Remaining bool
	Data      []byte
}

type IndexListener interface {
	Index(idx Index)
}

type EventListener interface {
	OnEvent(e Event)
	OnError(e Event, err error)
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
	OnListener(el EventListener)
	Listener() EventListener
	Topic(t string)
	OnTopic() string
	OnIndex(ix IndexListener)
	OId() int64
	OnOId(id int64)
}

type Index interface {
	Event
	Distributed() bool
}

type Postoffice interface {
	Send(e Event) error
	List(q Query)
	Recover(q Query)
}

type Publisher interface {
	Publish(e Event, ticket string) error
	Connect() error
	Close() error
}

type EventObj struct {
	Callback EventListener `json:"-"`
	core.PersistentableObj
	topic string `json:"-"`
	oid   int64  `json:"-"`
}

func (s *EventObj) Topic(t string) {
	s.topic = t
}

func (s *EventObj) OnTopic() string {
	return s.topic
}

func (s *EventObj) Distributed() bool {
	return false
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
	s.oid = oid
}

func (s *EventObj) OId() int64 {
	return s.oid
}
