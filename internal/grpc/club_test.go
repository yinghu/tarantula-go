package grpc

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestClub(t *testing.T) {
	p := Person{Name: "tester", Id: 100, Email: "tester@mail.com", Type: MemberType_Diamond, Age: 20}
	c := Club{Name: "abc", Id: 3}
	c.Members = append(c.Members, &p)
	out, err := proto.Marshal(&c)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	var cx Club
	err = proto.Unmarshal(out, &cx)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	if c.Name != "abc" && c.Name != cx.Name {
		t.Errorf("club name should be abc %s : %s", c.Name, cx.Name)
	}
	if c.Id != 3 && c.Id != cx.Id {
		t.Errorf("club ID should 3 be %d %d", c.Id, cx.Id)
	}
	if len(c.Members) != 1 {
		t.Errorf("member size should be 1 :%d", len(c.Members))
	}
	if len(cx.Members) != 1 {
		t.Errorf("member size should be 1 :%d", len(cx.Members))
	}
}
