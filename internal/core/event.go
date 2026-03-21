package core

import "time"

const (
	EVENT_FACTORY_ID uint32 = 1
)

type EventListener interface {
	OnEvent(e Event)
	OnError(e Event, err error)
}

type EventService interface {
	EventCreator
	VerifyTicket(ticket string) (OnSession, error)
	EventListener
	Publish(e Event) error
	List(q Query)
	Subscribe(topic string, listener EventListener) error
}

type Event interface {
	Inbound(buff DataBuffer) error
	Outbound(buff DataBuffer) error
	Persistentable
	OnListener(el EventListener)
	Listener() EventListener

	OnTopic(t string)
	Topic() string
	OId() int64
	OnOId(id int64)
	OnRecipientId(id int64)
	RecipientId() int64
}

type EventCreator interface {
	Create(classId uint32, topic string) (Event, error)
}

type Publisher interface {
	Publish(e Event, ticket string) error
	Connect() error
	Close() error
}

type Pusher interface {
	Push(e Event)
}

type Query interface {
	QId() uint32
	QFactoryId() uint32
	QClassId() uint32
	QTag() string
	QTopic() string
	QStartTime() time.Time
	QEndTime() time.Time
	QLimit() int32
	QOffset() int32
	QRead(b DataBuffer) error
	QWrite(b DataBuffer) error
	QCc() chan Chunk
	QEvent() Event
	QFilter(k, v []byte) bool
}

type EventObj struct {
	Callback EventListener `json:"-"`
	PersistentableObj
	ETopic string `json:"Topic"`
	EOid   int64  `json:"Oid,string"`
}

func (s *EventObj) OnTopic(t string) {
	s.ETopic = t
}

func (s *EventObj) Topic() string {
	return s.ETopic
}

func (s *EventObj) Inbound(buff DataBuffer) error {
	return nil
}
func (s *EventObj) Outbound(buff DataBuffer) error {
	return nil
}
func (s *EventObj) OnListener(el EventListener) {
	s.Callback = el
}
func (s *EventObj) Listener() EventListener {
	return s.Callback
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

func (s *EventObj) OnRecipientId(recipientId int64) {

}

func (s *EventObj) FactoryId() uint32 {
	return EVENT_FACTORY_ID
}
