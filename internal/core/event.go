package core

import (
	"time"

	"gameclustering.com/internal/protocol"
)

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
	//Publish(e Event) error
	Publish(e *protocol.Topic) error

	List(q Query)
	Subscribe(topic string, listener EventListener) error
	Unsubscribe(topic string) error
}

type Event interface {
	Inbound(buff DataBuffer) error
	Outbound(buff DataBuffer) error
	Persistentable
	OnListener(el EventListener)
	Listener() EventListener

	OnNodeId(t string)
	NodeId() string

	OnTag(t string)
	Tag() string

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
	ENodeId string `json:"nodeId"`
	ETag    string `json:"tag"`
	ETopic  string `json:"topic"`
	EOid    int64  `json:"oid,string"`
}

func (s *EventObj) OnNodeId(t string) {
	s.ENodeId = t
}

func (s *EventObj) NodeId() string {
	return s.ENodeId
}

func (s *EventObj) OnTag(t string) {
	s.ETag = t
}

func (s *EventObj) Tag() string {
	return s.ETag
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
