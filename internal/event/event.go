package event

import (
	"gameclustering.com/internal/persistence"
)

type Chunk struct {
	Remaining bool
	Data      []byte
}

type Event interface {
	Topic() bool
	Streaming(c Chunk) error
	persistence.Persistentable
}
