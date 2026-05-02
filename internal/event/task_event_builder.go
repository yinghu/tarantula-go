package event

import (
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskEventBuilder struct {
	Target *protocol.TaskEvent
	vb     *JobEventBuilder
	jb     *JobEventBuilder
}

func NewTaskEventBuilder(meta *protocol.Meta) *TaskEventBuilder {
	return &TaskEventBuilder{Target: &protocol.TaskEvent{Meta: meta}}
}

func (t *TaskEventBuilder) Start(ts *timestamppb.Timestamp) *TaskEventBuilder {
	t.Target.Start = ts
	return t
}

func (t *TaskEventBuilder) End(ts *timestamppb.Timestamp) *TaskEventBuilder {
	t.Target.End = ts
	return t
}

func (t *TaskEventBuilder) Description(desc string) *TaskEventBuilder {
	t.Target.Description = desc
	return t
}

func (t *TaskEventBuilder) ValidatorBuilder(meta *protocol.Meta) *JobEventBuilder {
	t.vb = NewJobEventBuilder(meta)
	return t.vb
}
func (t *TaskEventBuilder) JobBuilder(meta *protocol.Meta) *JobEventBuilder {
	t.jb = NewJobEventBuilder(meta)
	return t.jb
}

func (t *TaskEventBuilder) Build() *protocol.TaskEvent {
	t.Target.Validator = t.vb.Target
	t.Target.Job = t.jb.Target
	return t.Target
}
