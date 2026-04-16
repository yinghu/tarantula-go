package persistence

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/protocol"
)

func TestTaskBuilder(t *testing.T) {
	login := protocol.LoginObject{Name: "p100", Password: "passsword", SystemId: 100, AccessControl: 1, ReferenceId: 1}
	mf := NewLoginObjectFactory()
	kv, _ := mf.FromLoginObject(&login)
	tb := NewTaskBuilder(&protocol.Meta{NodeId: "node1", Tag: "presence", Name: "register", Id: 1, Prefix: 100})
	tb.Add(&protocol.Transaction{Meta: &protocol.Meta{NodeId: "node1", Tag: "inventory", Name: "grant", TaskId: 1, Id: 10, Prefix: 100}, Object: kv})
	tb.Add(&protocol.Transaction{Meta: &protocol.Meta{NodeId: "node1", Tag: "asset", Name: "update", TaskId: 1, Id: 20, Prefix: 200}, Object: kv})
	task := tb.Task()
	fmt.Printf("task meta %v\n", task.Meta)
	fmt.Printf("task trans 0 %v\n", task.Transactions[0])
	fmt.Printf("task trans 1 %v\n", task.Transactions[1])
	req, err := tb.Request()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	fmt.Printf("request %v\n", req)
	tbx := TaskBuilder{Target: task}
	reqx, err := tbx.Request()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	fmt.Printf("request %v\n", reqx)

	tbc := TaskBuilder{}
	tk, err := tbc.From(req.Data.Value)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	fmt.Printf("task meta %v\n", tk.Meta)
	fmt.Printf("task trans 0 %v\n", tk.Transactions[0])
	fmt.Printf("task trans 1 %v\n", tk.Transactions[1])

}
