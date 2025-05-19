package persistence

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/core"
)

type sample struct {
	Id        int64
	Name      string
	Age       int32
	Address   string
	Validated bool
	Pay       complex64
	core.PersistentableObj
}

func (s *sample) Write(value core.DataBuffer) error {
	value.WriteInt32(s.Age)
	value.WriteString(s.Name)
	value.WriteString(s.Address)
	value.WriteBool(s.Validated)
	value.WriteComplex64(s.Pay)
	return nil
}

func (s *sample) WriteKey(value core.DataBuffer) error {
	value.WriteInt64(s.Id)
	return nil
}

func (s *sample) Read(value core.DataBuffer) error {
	age, err := value.ReadInt32()
	if err != nil {
		return err
	}
	s.Age = age
	name, err := value.ReadString()
	if err != nil {
		return err
	}
	s.Name = name
	address, err := value.ReadString()
	if err != nil {
		return err
	}
	s.Address = address
	validated, err := value.ReadBool()
	if err != nil {
		return err
	}
	s.Validated = validated
	pay, err := value.ReadComplex64()
	if err != nil {
		return err
	}
	s.Pay = pay
	return nil
}

func TestLocalStore(t *testing.T) {
	local := LocalStore{InMemory: false, Path: "/home/yinghu/local", KeySize: 100, ValueSize: 200}
	err := local.Open()
	if err != nil {
		t.Errorf("Local store error %s", err.Error())
	}
	defer local.Close()
	sample1 := sample{Id: 200, Name: "yinghu12389", Address: "19809 150TH", Age: 9, Validated: true, Pay: 100}
	err = local.New(&sample1)
	if err!=nil{
		fmt.Printf("NO SAVE %s\n",err.Error())
	}
	load := sample{Id: 100}
	local.Load(&load)
	fmt.Printf("DATA :%v\n", load)

	load1 := sample{Id: 200}
	local.Load(&load1)
	fmt.Printf("DATA :%v\n", load1)
}
