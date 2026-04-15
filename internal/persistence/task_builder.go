package persistence

import "gameclustering.com/internal/protocol"

const (
	TASK_FACTORY_ID uint32 = 3
)

type TaskBuilder struct {
	target *protocol.Task
}

func NewTaskBuilder(meta *protocol.Meta) *TaskBuilder {
	trans := make([]*protocol.Transaction, 0)
	return &TaskBuilder{target: &protocol.Task{Meta: meta, Transactions: trans}}
}

func (b *TaskBuilder) Add(t *protocol.Transaction) *TaskBuilder {
	b.target.Transactions = append(b.target.Transactions, t)
	return b
}

//query task
func (b *TaskBuilder) Task() *protocol.Task {
	return b.target
}

//query request
func (b *TaskBuilder) Request() *protocol.Request {
	return &protocol.Request{}
}
