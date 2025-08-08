package persistence

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type sample struct {
	Tag  string
	Id   int64
	Str  string
	I32  int32
	B    bool
	C64  complex64
	C128 complex128
	I8   int8
	I16  int16
	I64  int64
	F32  float32
	F64  float64

	core.PersistentableObj
}

func (s *sample) ClassId() int {
	return 12
}
func (s *sample) ETag() string {
	return "sample:"
}
func (s *sample) Write(value core.DataBuffer) error {
	value.WriteInt64(s.I64)
	value.WriteInt32(s.I32)
	value.WriteInt16(s.I16)
	value.WriteInt8(s.I8)
	value.WriteString(s.Str)
	value.WriteBool(s.B)
	value.WriteComplex64(s.C64)
	value.WriteComplex128(s.C128)
	value.WriteFloat32(s.F32)
	value.WriteFloat64(s.F64)
	return nil
}

func (s *sample) WriteKey(value core.DataBuffer) error {
	value.WriteString(s.Tag)
	value.WriteInt64(s.Id)
	return nil
}

func (s *sample) ReadKey(value core.DataBuffer) error {
	tag, err := value.ReadString()
	if err != nil {
		return err
	}
	s.Tag = tag
	id, err := value.ReadInt64()
	if err != nil {
		return err
	}
	s.Id = id
	return nil
}

func (s *sample) Read(value core.DataBuffer) error {
	i64, err := value.ReadInt64()
	if err != nil {
		return err
	}
	s.I64 = i64
	i32, err := value.ReadInt32()
	if err != nil {
		return err
	}
	s.I32 = i32
	i16, err := value.ReadInt16()
	if err != nil {
		return err
	}
	s.I16 = i16
	i8, err := value.ReadInt8()
	if err != nil {
		return err
	}
	s.I8 = i8
	str, err := value.ReadString()
	if err != nil {
		return err
	}
	s.Str = str
	b, err := value.ReadBool()
	if err != nil {
		return err
	}
	s.B = b

	c64, err := value.ReadComplex64()
	if err != nil {
		return err
	}
	s.C64 = c64

	c128, err := value.ReadComplex128()
	if err != nil {
		return err
	}
	s.C128 = c128
	f32, err := value.ReadFloat32()
	if err != nil {
		return err
	}
	s.F32 = f32
	f64, err := value.ReadFloat64()
	if err != nil {
		return err
	}
	s.F64 = f64
	return nil
}

func TestDatatBuffer(t *testing.T) {
	local := BadgerLocal{InMemory: false, Path: "/home/yinghu/local/test"}
	err := local.Open()
	if err != nil {
		t.Errorf("Local store error %s", err.Error())
	}
	defer local.Close()
	sample1 := sample{Id: 200, Tag: "sample:"}
	sample1.I64 = 64
	sample1.I32 = 32
	sample1.I16 = 16
	sample1.I8 = 8
	sample1.Str = "hello"
	sample1.B = true
	sample1.C64 = 164
	sample1.C128 = 228
	sample1.F32 = 12.09
	sample1.F64 = 64.09
	local.Save(&sample1)

	load1 := sample{Id: 200, Tag: "sample:"}
	err = local.Load(&load1)
	if err != nil {
		t.Errorf("no load %s", err.Error())
	}
	if load1.C64 != 164 {
		t.Errorf("no load %s", "c64")
	}
	if load1.C128 != 228 {
		t.Errorf("no load %s", "c128")
	}
	if !load1.B {
		t.Errorf("no load %s", "bool")
	}
	if load1.Str != "hello" {
		t.Errorf("no load %s", "string")
	}
	if load1.I64 != 64 {
		t.Errorf("no load %s", "i64")
	}

	if load1.I32 != 32 {
		t.Errorf("no load %s", "i32")
	}
	if load1.I16 != 16 {
		t.Errorf("no load %s", "i16")
	}
	if load1.I8 != 8 {
		t.Errorf("no load %s", "i8")
	}
	if load1.F32 != 12.09 {
		t.Errorf("no load %f", load1.F32)
	}
	if load1.F64 != 64.09 {
		t.Errorf("no load %f", load1.F64)
	}
	fmt.Printf("Rev :%d\n", load1.Revision())
	sample2 := sample{Id: 1000, Tag: "sample:"}
	sample2.Str = "test"
	local.Save(&sample2)

	sample3 := sample{Id: 2000, Tag: "sample:"}
	sample3.Str = "high"
	local.Save(&sample3)

	sample4 := sample{Id: 20, Tag: "sample:"}
	sample4.Str = "lower"
	local.Save(&sample4)

	px := BufferProxy{}
	px.NewProxy(100)
	px.WriteString("sample:")
	px.Flip()
	ct := 0
	local.List(&px, func(k, v core.DataBuffer, rev uint64) bool {
		d := sample{}
		d.ReadKey(k)
		cid, _ := v.ReadInt32()
		fmt.Printf("header %d %d\n", cid, rev)
		d.Read(v)
		d.OnRevision(rev)
		fmt.Printf("streaming %s, %d\n", d.Tag, d.Id)
		fmt.Printf("streaming %s, %d\n", d.Str, d.I64)
		ct++
		return true
	})
	if ct != 4 {
		t.Errorf("should be 3 items %s", err.Error())
	}
	se := event.StatEvent{Tag: sample1.ETag(), Name: "total"}
	err = local.Load(&se)
	if err != nil {
		fmt.Printf("se load err %s\n", err.Error())
	}
	fmt.Printf("Total count : %d\n", se.Count)
}
