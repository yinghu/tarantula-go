package core

type PersistentableFactory interface {
	Create(classId int) Persistentable
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

type Stream func(k, v DataBuffer) bool

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

type DataStoreFactory interface {
	Create(name string) (DataStore, error)
}

type DataStore interface {
	Load(p Persistentable) error
	Save(p Persistentable) error
	List(prefix DataBuffer, s Stream) error
	Close() error
}
