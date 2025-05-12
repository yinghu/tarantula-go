package persistence

import (
	"errors"
	"fmt"

	buffer "github.com/0xc0d/encoding/bytebuffer"
	badger "github.com/dgraph-io/badger/v4"
)

type LocalStore struct {
	InMemory bool
	Path     string
	Db       *badger.DB
}

func (s *LocalStore) Set(key buffer.ByteBuffer, value buffer.ByteBuffer) error {
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

func (s *LocalStore) Get(key buffer.ByteBuffer) (buffer.ByteBuffer, error) {
	if key.Remaining() == 0 {
		return key, errors.New("bad key/value")
	}
	s.Db.View(func(txn *badger.Txn) error {
		k := make([]byte, key.Remaining())
		key.GetBytes(k, 0, key.Remaining())
		item, err := txn.Get(k)
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {
			fmt.Printf("GET : %s\n", val)
			return nil
		})
		return nil
	})
	return key, nil
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
