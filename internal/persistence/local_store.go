package persistence

import (
	badger "github.com/dgraph-io/badger/v4"
)

type LocalStore struct {
	InMemory bool
	Path     string
	Db       *badger.DB
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
