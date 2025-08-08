package persistence

import (
	"errors"
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"

	badger "github.com/dgraph-io/badger/v4"
)

const (
	BDG_KEY_SIZE   int = 200
	BDG_VALUE_SIZE int = 4096
)

type BadgerLocal struct {
	InMemory bool
	Path     string
	Db       *badger.DB
	Seq      core.Sequence
}

func (s *BadgerLocal) Save(t core.Persistentable) error {
	key := BufferProxy{}
	key.NewProxy(BDG_KEY_SIZE)
	value := BufferProxy{}
	value.NewProxy(BDG_VALUE_SIZE)
	t.WriteKey(&key)
	value.WriteInt32(int32(t.ClassId()))
	t.Write(&value)
	key.Flip()
	value.Flip()
	return s.set(&key, &value, t)
}

func (s *BadgerLocal) Load(t core.Persistentable) error {
	key := BufferProxy{}
	key.NewProxy(BDG_KEY_SIZE)
	t.WriteKey(&key)
	key.Flip()
	value := BufferProxy{}
	value.NewProxy(BDG_VALUE_SIZE)
	rev, err := s.get(&key, &value)
	if err != nil {
		return err
	}
	value.ReadInt32()
	t.Read(&value)
	t.OnRevision(rev)
	return nil
}

func (s *BadgerLocal) List(prefix core.DataBuffer, stream core.Stream) error {
	return s.Db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		p, err := prefix.Read(0)
		if err != nil {
			return err
		}
		key := BufferProxy{}
		key.NewProxy(BDG_KEY_SIZE)
		value := BufferProxy{}
		value.NewProxy(BDG_VALUE_SIZE)
		for it.Seek(p); it.ValidForPrefix(p); it.Next() {
			kv := it.Item()
			key.Clear()
			err = key.Write(kv.Key())
			if err != nil {
				return err
			}
			err = kv.Value(func(val []byte) error {
				value.Clear()
				return value.Write(val)
			})
			if err != nil {
				return err
			}
			key.Flip()
			value.Flip()
			if !stream(&key, &value, kv.Version()) {
				break
			}
		}
		return nil
	})
}

func (s *BadgerLocal) set(key *BufferProxy, value *BufferProxy, t core.Persistentable) error {
	if key.Remaining() == 0 || value.Remaining() == 0 {
		return errors.New("bad key/value")
	}
	return s.Db.Update(func(txn *badger.Txn) error {
		k, _ := key.Read(0)
		v, _ := value.Read(0)
		item, err := txn.Get(k)
		if err == nil && t.Revision() != item.Version() {
			return fmt.Errorf("rev not match %d %d", t.Revision(), item.Version())
		}
		if err == nil && t.Revision() == item.Version() {
			return txn.Set(k, v)
		}
		//new entry
		if err = txn.Set(k, v); err != nil {
			return err
		}
		//update stat total
		se := event.StatEvent{Tag: t.ETag(), Name: "total"}
		ckey := BufferProxy{}
		ckey.NewProxy(BDG_KEY_SIZE)
		cvalue := BufferProxy{}
		cvalue.NewProxy(BDG_VALUE_SIZE)
		se.WriteKey(&ckey)
		ckey.Flip()
		ck, _ := ckey.Read(0)
		citem, err := txn.Get(ck)
		if err == nil {
			err = citem.Value(func(val []byte) error {
				return cvalue.Write(val)
			})
			if err != nil {
				return err //rollback
			}
			cvalue.Flip()
			cvalue.ReadInt32()
			se.Read(&cvalue)
			se.Count = se.Count + 1
		} else {
			se.Count = 1
		}
		cvalue.Clear()
		cvalue.WriteInt32(int32(se.ClassId()))
		se.Write(&cvalue)
		cvalue.Flip()
		cv, _ := cvalue.Read(0)
		return txn.Set(ck, cv)

	})
}

func (s *BadgerLocal) get(key *BufferProxy, value *BufferProxy) (uint64, error) {
	if key.Remaining() == 0 {
		return 0, errors.New("bad key/value")
	}
	var rev uint64
	err := s.Db.View(func(txn *badger.Txn) error {
		k, _ := key.Read(0)
		item, err := txn.Get(k)
		if err != nil {
			return err
		}
		rev = item.Version()
		item.Value(func(val []byte) error {
			value.Write(val)
			value.Flip()
			return nil
		})
		return nil
	})
	if err != nil {
		return 0, err
	}
	return rev, nil
}

func (s *BadgerLocal) Open() error {
	var opt badger.Options
	if s.InMemory {
		opt = badger.DefaultOptions("").WithInMemory(true)
	} else {
		opt = badger.DefaultOptions(s.Path)
		opt.SyncWrites = false
	}
	db, err := badger.Open(opt)
	if err != nil {
		return err
	}
	s.Db = db
	return nil
}

func (s *BadgerLocal) Close() error {
	s.Db.Sync()
	return s.Db.Close()
}

func (s *BadgerLocal) GC() error {
	return s.Db.RunValueLogGC(0.7)
}
