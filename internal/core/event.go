package core

import (
	"time"
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
	QNodeId() string
	QTag() string
	QTopic() string
	QStartTime() time.Time
	QEndTime() time.Time
	QLimit() int32
	QOffset() int32
	QRead(b DataBuffer) error
	QWrite(b DataBuffer) error
	QEvent() Event
	QFilter(k, v []byte) bool
}

type QueryFactory interface {
	Export(query Query) ([]byte, error)
	Import(criteria []byte) (Query, error)
	Query() Query
}
