package persistence

type PersistentableFactory interface {
	Create(classId int) Persistentable
}

type DataBuffer interface {
	WriteBool(data bool)

	WriteComplex64(data complex64)

	WriteComplex128(data complex128)

	WriteFloat64(data float64)

	WriteFloat32(data float32)

	WriteInt64(data int64)

	WriteInt32(data int32)

	WriteInt16(data int16)

	WriteInt8(data int8)

	WriteString(data string)

	ReadInt32() int32

	ReadInt64() int64

	ReadFloat32() float32

	ReadFloat64() float64

	ReadInt16() int16

	ReadInt8() int8

	ReadComplex64() complex64

	ReadComplex128() complex128

	ReadString() string

	ReadBool() bool
	Read() ([]byte, error)

	Write(data []byte)
	Remaining() int

	Flip() error
}

type Persistentable interface {
	Write(value DataBuffer) error
	WriteKey(key DataBuffer) error
	Read(value DataBuffer) error
	ReadKey(key DataBuffer) error
	ClassId() int
	Revision() int64
}

type PersistentableObj struct {
	Fid int
	Cid int
	Rev int64
}

func (s *PersistentableObj) Write(value DataBuffer) error {
	return nil
}

func (s *PersistentableObj) WriteKey(value DataBuffer) error {
	return nil
}

func (s *PersistentableObj) Read(value DataBuffer) error {
	return nil
}

func (s *PersistentableObj) ReadKey(value DataBuffer) error {
	return nil
}

func (s *PersistentableObj) ClassId() int {
	return s.Cid
}

func (s *PersistentableObj) Revision() int64 {
	return s.Rev
}
