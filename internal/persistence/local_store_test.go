package persistence

import (
	"fmt"
	"testing"

	buffer "github.com/0xc0d/encoding/bytebuffer"
)

type Sample struct {
	Id      uint64
	Name    string
	Age     uint32
	Address string
}

func (s *Sample) Write(value *buffer.ByteBuffer) error {
	lz := len(s.Name)
	value.PutUint32(uint32(lz))
	e := value.PutBytes([]byte(s.Name), 0, lz)
	if e != nil {
		fmt.Printf("Error:%s\n", e.Error())
	}
	fmt.Printf("POS: %d\n", value.Position())
	//value.PutUint32(s.Age)
	//lz1 := len(s.Address)
	//value.PutUint32(uint32(lz1))
	//value.PutBytes([]byte(s.Address), 0, lz1)
	return nil
}

func (s *Sample) WriteKey(value *buffer.ByteBuffer) error {
	value.PutUint64(s.Id)
	return nil
}

func (s *Sample) Read(value *buffer.ByteBuffer) error {
	fmt.Printf("RM %d %d\n", value.Remaining(), value.Position())
	len, e := value.GetAsUint32()
	if e != nil {
		fmt.Printf("Error:%s\n", e.Error())
	}
	fmt.Printf("LEN :%d\n", len)
	name := make([]byte, int(len))
	value.GetBytes(name, 0, int(len))
	s.Name = string(name)
	//age, _ := value.GetAsUint32()
	//s.Age = age
	//len1, _ := value.GetAsInt32()
	//address := make([]byte, int(len1))
	//value.GetBytes(address, 0, int(len1))
	//s.Address = string(address)
	return nil
}

func (s *Sample) ReadKey(value *buffer.ByteBuffer) error {
	return nil
}

func TestLocalStore(t *testing.T) {
	local := LocalStore{InMemory: true, Path: "/home/yinghu/local"}
	err := local.Open()
	if err != nil {
		t.Errorf("Local store error %s", err.Error())
	}
	defer local.Close()
	key := buffer.NewByteBuffer(10)
	key.PutUint16(10)
	key.Flip()
	value := buffer.NewByteBuffer(10)
	value.PutBytes([]byte("hello"), 0, 5)
	value.Flip()
	e := local.Set(key, value)
	if e != nil {
		t.Errorf("Local store set error %s", err.Error())
	}
	key.Rewind()
	ret, err := local.Get(key)
	if err != nil {
		t.Errorf("Local store get error %s", err.Error())
	}
	fmt.Printf("Buffer : %d\n", ret.Remaining())
	sample := Sample{Id: 100, Name: "test", Age: 12, Address: "1980 NE 159"}
	if local.Save(&sample) != nil {
		t.Errorf("Local store save error %s", err.Error())
	}
	load := Sample{Id: 100}
	if local.Load(&load) != nil {
		t.Errorf("Local store load error %s", err.Error())
	}
	fmt.Printf("LOADED %v\n", load)
}
