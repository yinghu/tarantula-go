package event

import (
	"gameclustering.com/internal/persistence"
)

type Event interface {
	Topic() bool
	Send() error
	persistence.Persistentable
}
