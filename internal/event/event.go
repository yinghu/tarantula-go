package event

type Event interface {
	OnTopic() bool
}


