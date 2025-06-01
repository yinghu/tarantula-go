package cluster

import (
	"iter"
)

type Ctx interface {
	Put(key string, value string) error
	Get(key string) (string, error)
	Del(key string) error
}

type Exec func(ctx Ctx)

type Cluster interface {
	Local() Node
	View() iter.Seq[Node]
	Partition(key []byte) Node
	Atomic(t Exec) error
}
