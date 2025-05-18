package persistence

import (
	"errors"

	"gameclustering.com/internal/core"

	badger "github.com/dgraph-io/badger/v4"
)

type LocalStore struct {
	InMemory  bool
	Path      string
	Db        *badger.DB
	KeySize   int
	ValueSize int
}

func (s *LocalStore) Save(t core.Persistentable) error {

	key := BufferProxy{}
	key.NewProxy(s.KeySize)
	value := BufferProxy{}
	value.NewProxy(s.ValueSize)
	t.WriteKey(&key)
	t.Write(&value)
	key.Flip()
	value.Flip()
	return s.Set(&key, &value)
}

func (s *LocalStore) Load(t core.Persistentable) error {
	key := BufferProxy{}
	key.NewProxy(s.KeySize)
	t.WriteKey(&key)
	key.Flip()
	value := BufferProxy{}
	value.NewProxy(s.ValueSize)
	err := s.Get(&key, &value)
	if err != nil {
		return err
	}
	t.Read(&value)
	return nil
}

func (s *LocalStore) Set(key *BufferProxy, value *BufferProxy) error {
	if key.Remaining() == 0 || value.Remaining() == 0 {
		return errors.New("bad key/value")
	}

	return s.Db.Update(func(txn *badger.Txn) error {
		k, _ := key.Read(0)
		v, _ := value.Read(0)
		return txn.Set(k, v)
	})
}

func (s *LocalStore) Get(key *BufferProxy, value *BufferProxy) error {
	if key.Remaining() == 0 {
		return errors.New("bad key/value")
	}
	err := s.Db.View(func(txn *badger.Txn) error {
		k, _ := key.Read(0)
		item, err := txn.Get(k)
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {
			value.Write(val)
			value.Flip()
			return nil
		})
		return nil
	})
	if err != nil {
		return err
	}
	return nil
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
