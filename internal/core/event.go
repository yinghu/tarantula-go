package core

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
