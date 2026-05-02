package event

import (
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type JobEventBuilder struct {
	Target *protocol.JobEvent
	tb     *TransactionEventBuilder
}

func NewJobEventBuilder(meta *protocol.Meta) *JobEventBuilder {
	return &JobEventBuilder{Target: &protocol.JobEvent{Meta: meta, Transactions: make([]*protocol.TransactionEvent, 0)}}
}

func (t *JobEventBuilder) Start(ts *timestamppb.Timestamp) *JobEventBuilder {
	t.Target.Start = ts
	return t
}

func (t *JobEventBuilder) End(ts *timestamppb.Timestamp) *JobEventBuilder {
	t.Target.End = ts
	return t
}

func (t *JobEventBuilder) Description(desc string) *JobEventBuilder {
	t.Target.Description = desc
	return t
}

// chaining build a transaction New to Build
func (t *JobEventBuilder) Transaction() *TransactionEventBuilder {
	t.tb = NewTransactionEventBuilder(t)
	return t.tb
}
