package persistence

type PersistentableFactory interface {
	Create(classId int) Persistentable
}

type Persistentable interface {
	Write(value *BufferProxy) error
	WriteKey(key *BufferProxy) error
	Read(value *BufferProxy) error
	ReadKey(key *BufferProxy) error
	ClassId() int
	Revision() int64
}

type PersistentableObj struct {
	Fid int
	Cid int
	Rev int64
}

func (s *PersistentableObj) Write(value *BufferProxy) error {
	return nil
}

func (s *PersistentableObj) WriteKey(value *BufferProxy) error {
	return nil
}

func (s *PersistentableObj) Read(value *BufferProxy) error {
	return nil
}

func (s *PersistentableObj) ReadKey(value *BufferProxy) error {
	return nil
}

func (s *PersistentableObj) ClassId() int {
	return s.Cid
}

func (s *PersistentableObj) Revision() int64 {
	return s.Rev
}
