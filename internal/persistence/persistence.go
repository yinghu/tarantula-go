package persistence

type Persistentable interface {
	Write(value *BufferProxy) error
	WriteKey(key *BufferProxy) error
	Read(value *BufferProxy) error
	ReadKey(key *BufferProxy) error
}
