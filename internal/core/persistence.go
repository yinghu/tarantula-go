package core


type Persistentable interface {
	Write(value DataBuffer) error
	WriteKey(key DataBuffer) error
	Read(value DataBuffer) error
	ReadKey(key DataBuffer) error
	ClassId() int
	Revision() int64
	Timestamp() int64
	OnTimestamp(tsp int64)
	OnRevision(rev int64)
	ETag() string
}

type PersistentableObj struct {
	Rev int64 `json:"rev,string"`
	Tsp int64  `json:"timestamp,string"`
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
	return 0
}

func (s *PersistentableObj) Revision() int64 {
	return s.Rev
}
func (s *PersistentableObj) Timestamp() int64 {
	return s.Tsp
}
func (s *PersistentableObj) OnTimestamp(tsp int64) {
	s.Tsp = tsp
}

func (s *PersistentableObj) OnRevision(rev int64) {
	s.Rev = rev
}

func (s *PersistentableObj) ETag() string {
	return "ent:"
}

type DataStore interface {
	Load(p Persistentable) error
	Save(p Persistentable) error
	List(prefix DataBuffer, s Stream) error
	Close() error
}
