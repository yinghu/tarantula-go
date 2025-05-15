package persistence

import (
	"fmt"
	"testing"
)

type Sample struct {
	Id        int64
	Name      string
	Age       int32
	Address   string
	Validated bool
	Pay       complex64
}

func (s *Sample) Write(value *BufferProxy) error {
	value.WriteInt32(s.Age)
	value.WriteString(s.Name)
	value.WriteString(s.Address)
	value.WriteBool(s.Validated)
	value.WriteComplex64(s.Pay)
	return nil
}

func (s *Sample) WriteKey(value *BufferProxy) error {
	value.WriteInt64(s.Id)
	return nil
}

func (s *Sample) Read(value *BufferProxy) error {
	s.Age = value.ReadInt32()
	s.Name = value.ReadString()
	s.Address = value.ReadString()
	s.Validated = value.ReadBool()
	s.Pay = value.ReadComplex64()
	return nil
}

func (s *Sample) ReadKey(value *BufferProxy) error {
	s.Id = value.ReadInt64()
	return nil
}

func TestLocalStore(t *testing.T) {
	local := LocalStore{InMemory: false, Path: "/home/yinghu/local", KeySize: 100, ValueSize: 200}
	err := local.Open()
	if err != nil {
		t.Errorf("Local store error %s", err.Error())
	}
	defer local.Close()
	sample := Sample{Id: 200, Name: "yinghu", Address: "19809 150TH", Age: 9, Validated: true, Pay: 100}
	local.Save(&sample)
	load := Sample{Id: 100}
	local.Load(&load)
	fmt.Printf("DATA :%v\n", load)

	load1 := Sample{Id: 200}
	local.Load(&load1)
	fmt.Printf("DATA :%v\n", load1)
}
