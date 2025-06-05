package persistence

import (
	"errors"

	"gameclustering.com/internal/core"

	badger "github.com/dgraph-io/badger/v4"
)

const (
	BDG_KEY_SIZE   int = 200
	BDG_VALUE_SIZE int = 1800
)

type Cache struct {
	InMemory  bool
	Path      string
	Db        *badger.DB
	Seq       core.Sequence
}

func (s *Cache) Save(t core.Persistentable) error {
	key := BufferProxy{}
	key.NewProxy(BDG_KEY_SIZE)
	value := BufferProxy{}
	value.NewProxy(BDG_VALUE_SIZE)
	t.WriteKey(&key)
	t.Write(&value)
	key.Flip()
	value.Flip()
	return s.Set(&key, &value)
}

func (s *Cache) New(t core.Persistentable) error {
	key := BufferProxy{}
	key.NewProxy(BDG_KEY_SIZE)
	value := BufferProxy{}
	value.NewProxy(BDG_VALUE_SIZE)
	t.WriteKey(&key)
	t.Write(&value)
	key.Flip()
	value.Flip()
	return s.SetNew(&key, &value)
}

func (s *Cache) Load(t core.Persistentable) error {
	key := BufferProxy{}
	key.NewProxy(BDG_KEY_SIZE)
	t.WriteKey(&key)
	key.Flip()
	value := BufferProxy{}
	value.NewProxy(BDG_VALUE_SIZE)
	err := s.Get(&key, &value)
	if err != nil {
		return err
	}
	t.Read(&value)
	return nil
}

func (s *Cache) SetNew(key *BufferProxy, value *BufferProxy) error {
	if key.Remaining() == 0 || value.Remaining() == 0 {
		return errors.New("bad key/value")
	}

	return s.Db.Update(func(txn *badger.Txn) error {
		k, _ := key.Read(0)
		v, _ := value.Read(0)
		e := badger.NewEntry(k, v)
		return txn.SetEntry(e)
	})
}

func (s *Cache) Set(key *BufferProxy, value *BufferProxy) error {
	if key.Remaining() == 0 || value.Remaining() == 0 {
		return errors.New("bad key/value")
	}

	return s.Db.Update(func(txn *badger.Txn) error {
		k, _ := key.Read(0)
		v, _ := value.Read(0)
		return txn.Set(k, v)
	})
}

func (s *Cache) Get(key *BufferProxy, value *BufferProxy) error {
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

func (s *Cache) Open() error {
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

func (s *Cache) Close() error {
	s.Db.Sync()
	return s.Db.Close()

}
