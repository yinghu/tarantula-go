package persistence

type Persistentable interface {
	Write(value *BufferProxy) error
	WriteKey(key *BufferProxy) error
	Read(value *BufferProxy) error
	ReadKey(key *BufferProxy) error
}

type PersistentableObj struct {
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
