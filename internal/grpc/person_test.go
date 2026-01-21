package grpc

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestPerson(t *testing.T) {
	p := Person{Name: "tester", Id: 100, Email: "tester@mail.com", Type: MemberType_Diamond, Age: 20}
	if p.Name != "tester" {
		t.Errorf("person name should be tester %s", p.Name)
	}
	if p.Email != "tester@mail.com" {
		t.Errorf("person email should be tester@mail.com %s", p.Email)
	}
	if p.Type != MemberType_Diamond {
		t.Errorf("person type should be %d %d", MemberType_Diamond, p.Type)
	}
	if p.Id != 100 {
		t.Errorf("person ID should be %d %d", 100, p.Id)
	}
	if p.Age != 20 {
		t.Errorf("person age should be %d %d", 20, p.Age)
	}
	out, err := proto.Marshal(&p)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	var px Person
	err = proto.Unmarshal(out,&px)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	if p.Name != px.Name{
		t.Errorf("person names should be same %s : %s",p.Name,px.Name)
	}
	if p.Email != px.Email{
		t.Errorf("person email should be same %s : %s",p.Email,px.Email)
	}
	if p.Type != px.Type{
		t.Errorf("person type should be same %d : %d",p.Type,px.Type)
	}
	if p.Id != px.Id{
		t.Errorf("person ID should be same %d : %d",p.Id,px.Id)
	}
	if p.Age != px.Age{
		t.Errorf("person age should be same %d : %d",p.Age,px.Age)
	}
}
