package persistence

import (
	"fmt"
	"testing"
)

type Sample struct {
	Id      int32
	Name    string
	Age     int32
	Address string
}

func (s *Sample) Write(value *BufferProxy) error {
	value.WriteInt32(s.Age)
	value.WriteString(s.Name)
	value.WriteString(s.Address)
	return nil
}

func (s *Sample) WriteKey(value *BufferProxy) error {
	value.WriteInt32(s.Id)
	return nil
}

func (s *Sample) Read(value *BufferProxy) error {
	s.Age = value.ReadInt32()
	s.Name = value.ReadString()
	s.Address = value.ReadString()
	return nil
}

func (s *Sample) ReadKey(value *BufferProxy) error {
	return nil
}

func TestLocalStore(t *testing.T) {
	local := LocalStore{InMemory: true, Path: "/home/yinghu/local"}
	err := local.Open()
	if err != nil {
		t.Errorf("Local store error %s", err.Error())
	}
	defer local.Close()
	sample := Sample{Id: 100, Name: "teter", Address: "1980 150TH",Age: 9}
	local.Save(&sample)
	load := Sample{Id: 100}
	local.Load(&load)
	fmt.Printf("DATA :%v\n",load)
}
