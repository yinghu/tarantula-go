package persistence

type PersistentableFactory interface {
	Create(classId int) Persistentable
}

type DataBuffer interface {
	WriteBool(data bool) error

	WriteComplex64(data complex64) error

	WriteComplex128(data complex128) error

	WriteFloat64(data float64) error

	WriteFloat32(data float32) error

	WriteInt64(data int64) error

	WriteInt32(data int32) error

	WriteInt16(data int16) error

	WriteInt8(data int8) error

	WriteString(data string) error

	ReadInt32() (int32, error)

	ReadInt64() (int64, error)

	ReadFloat32() (float32, error)

	ReadFloat64() (float64, error)

	ReadInt16() (int16, error)

	ReadInt8() (int8, error)

	ReadComplex64() (complex64, error)

	ReadComplex128() (complex128, error)

	ReadString() (string, error)

	ReadBool() (bool, error)
	Read() ([]byte, error)

	Write(data []byte) error
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
