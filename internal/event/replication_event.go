package event

import (
	"gameclustering.com/internal/core"
)

type ReplicationEvent struct {
	Key   []byte
	Value []byte
	core.PersistentableObj
}
