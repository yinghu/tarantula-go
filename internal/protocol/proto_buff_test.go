package protocol

import (
	"fmt"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestEvent(t *testing.T) {
	obj, err := anypb.New(&MessageEvent{Title: "test", Message: "message", Source: "any"})
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
		return
	}
	e := Event{Header: &Header{FactoryId: 1, ClassId: 1}, Message: obj}
	data, err := proto.Marshal(&e)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
		return
	}
	fmt.Printf("size of data %d\n", len(data))

	var ex Event
	err = proto.Unmarshal(data, &ex)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
		return
	}
	fmt.Printf("header %v\n", ex.Header)
	var me MessageEvent
	err = ex.Message.UnmarshalTo(&me)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
		return
	}
	fmt.Printf("data %v\n", &me)
}
