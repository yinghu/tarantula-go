package cluster

import (
	"iter"
)

type Cluster interface {
	Local() Node
	View() iter.Seq[Node]
	Partition(key []byte) Node
}
