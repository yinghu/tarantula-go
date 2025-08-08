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
	Revision() uint64
	OnRevision(rev uint64)
	ETag() string
}

type PersistentableObj struct {
	Fid int
	Cid int
	Rev uint64
}

type Stream func(k, v DataBuffer,rev uint64) bool

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

func (s *PersistentableObj) Revision() uint64 {
	return s.Rev
}

func (s *PersistentableObj) OnRevision(rev uint64) {
	s.Rev = rev
}

func (s *PersistentableObj) ETag() string {
	return "ent:"
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
