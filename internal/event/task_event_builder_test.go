package event

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTaskEventBuilder(t *testing.T) {
	tb := NewTaskEventBuilder(&protocol.Meta{Id: 1, Name: "task"}).Description("task")
	vb := tb.ValidatorBuilder(&protocol.Meta{Id: 2, Name: "validaor"})
	vb.Transaction().New(&protocol.Meta{Name: "v1"}).Build().New(&protocol.Meta{Name: "v2"}).Build()
	jb := tb.JobBuilder(&protocol.Meta{Id: 4, Name: "job"}).Description("desc")
	jb.Transaction().New(&protocol.Meta{Name: "j1"}).Build().New(&protocol.Meta{Name: "J2"}).Build()
	jb.Start(timestamppb.Now()).End(timestamppb.Now())
	vb.Description("validator")
	te := tb.Build()
	if len(te.Validator.Transactions) != 2 {
		t.Errorf("validator should have 2 transactions %d", len(te.Validator.Transactions))
	}
	if len(te.Job.Transactions) != 2 {
		t.Errorf("job should have 2 transactions %d", len(te.Job.Transactions))
	}
	fmt.Printf("Task event %v", te)
}
