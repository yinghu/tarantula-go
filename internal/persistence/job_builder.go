package persistence

import "gameclustering.com/internal/protocol"

type JobBuilder struct {
	Target *protocol.Job
}

func NewJobBuilder(meta *protocol.Meta) *JobBuilder {
	trans := make([]*protocol.Transaction, 0)
	return &JobBuilder{Target: &protocol.Job{Meta: meta, Transactions: trans}}
}

func (jb *JobBuilder) Add(t *protocol.Transaction) *JobBuilder {
	jb.Target.Transactions = append(jb.Target.Transactions, t)
	return jb
}

func (jb *JobBuilder) Job() *protocol.Job {
	return jb.Target
}
