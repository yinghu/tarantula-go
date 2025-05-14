package persistence

import (
	"errors"
	"fmt"

	buffer "github.com/0xc0d/encoding/bytebuffer"
	badger "github.com/dgraph-io/badger/v4"
)

type Persistentable interface {
	Write(value *buffer.ByteBuffer) error
	WriteKey(key *buffer.ByteBuffer) error
	Read(value *buffer.ByteBuffer) error
	ReadKey(key *buffer.ByteBuffer) error
}

type LocalStore struct {
	InMemory bool
	Path     string
	Db       *badger.DB
}

func (s *LocalStore) Save(t Persistentable) error {
	key := buffer.NewByteBuffer(100)
	value := buffer.NewByteBuffer(200)
	t.WriteKey(key)
	t.Write(value)
	key.Flip()
	value.Flip()
	fmt.Printf("RMV :%d\n", value.Remaining())
	return s.Set(key, value)
}

func (s *LocalStore) Load(t Persistentable) error {
	key := buffer.NewByteBuffer(100)
	t.WriteKey(key)
	key.Flip()
	value, err := s.Get(key)
	if err != nil {
		return err
	}
	t.Read(value)
	return nil
}

func (s *LocalStore) Set(key *buffer.ByteBuffer, value *buffer.ByteBuffer) error {
	if key.Remaining() == 0 || value.Remaining() == 0 {
		return errors.New("bad key/value")
	}

	return s.Db.Update(func(txn *badger.Txn) error {
		k := make([]byte, key.Remaining())
		v := make([]byte, value.Remaining())
		key.GetBytes(k, 0, key.Remaining())
		value.GetBytes(v, 0, value.Remaining())
		return txn.Set(k, v)
	})
}

func (s *LocalStore) Get(key *buffer.ByteBuffer) (*buffer.ByteBuffer, error) {
	if key.Remaining() == 0 {
		return key, errors.New("bad key/value")
	}
	value := buffer.NewByteBuffer(200)
	err := s.Db.View(func(txn *badger.Txn) error {
		k := make([]byte, key.Remaining())
		key.GetBytes(k, 0, key.Remaining())
		item, err := txn.Get(k)
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {
			fmt.Printf("GET : %d\n", len(val))
			value.PutBytes(val, 0, len(val))
			value.Flip()
			return nil
		})
		return nil
	})
	if err != nil {
		return key, err
	}
	return value, nil
}

func (s *LocalStore) Open() error {
	var opt badger.Options
	if s.InMemory {
		opt = badger.DefaultOptions("").WithInMemory(true)
	} else {
		opt = badger.DefaultOptions(s.Path)
	}
	db, err := badger.Open(opt)
	if err != nil {
		return err
	}
	s.Db = db
	return nil
}

func (s *LocalStore) Close() {
	s.Db.Close()
}
