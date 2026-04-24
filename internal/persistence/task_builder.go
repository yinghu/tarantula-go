package persistence

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	TASK_FACTORY_ID uint32 = 3
	TASK_CLASS_ID   uint32 = 1
)

type TaskBuilder struct {
	Target *protocol.Task
}

func NewTaskBuilder(meta *protocol.Meta) *TaskBuilder {
	jobs := make([]*protocol.Job, 0)
	return &TaskBuilder{Target: &protocol.Task{Meta: meta, Jobs: jobs}}
}

func (b *TaskBuilder) Add(j *protocol.Job) *TaskBuilder {
	b.Target.Jobs = append(b.Target.Jobs, j)
	return b
}

// query task
func (b *TaskBuilder) Task() *protocol.Task {
	return b.Target
}

// query request
func (b *TaskBuilder) Request() (*protocol.Request, error) {
	req := protocol.Request{Opt: core.CREATE_DATA_REQUEST, Data: &protocol.Data{}}
	if b.Target.Meta.Id <= 0 {
		return &req, fmt.Errorf("task id should be more than zero uint64")
	}
	buff := core.NewBuffer(8)
	buff.WriteUInt64(b.Target.Meta.Id)
	buff.Flip()
	key, err := buff.Read(0)
	if err != nil {
		return &req, err
	}
	value, err := proto.Marshal(b.Target)
	if err != nil {
		return &req, err
	}
	req.Data.Header = &protocol.Header{FactoryId: TASK_FACTORY_ID, ClassId: TASK_CLASS_ID, Mutable: true}
	req.Data.Key = key
	req.Data.Value = value
	return &req, nil
}

// parse
func (b *TaskBuilder) From(data []byte) (*protocol.Task, error) {
	b.Target = &protocol.Task{}
	err := proto.Unmarshal(data, b.Target)
	return b.Target, err
}
