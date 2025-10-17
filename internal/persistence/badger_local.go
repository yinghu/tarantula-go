package persistence

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"

	badger "github.com/dgraph-io/badger/v4"
)

const (
	BDG_KEY_SIZE   int = 200
	BDG_VALUE_SIZE int = 4096
)

type BadgerLocal struct {
	InMemory    bool
	Path        string
	Db          *badger.DB
	LogDisabled bool
}

func (s *BadgerLocal) Save(t core.Persistentable) error {
	key := BufferProxy{}
	key.NewProxy(BDG_KEY_SIZE)
	value := BufferProxy{}
	value.NewProxy(BDG_VALUE_SIZE)
	t.WriteKey(&key)
	value.WriteInt32(int32(t.ClassId()))
	value.WriteInt64(int64(t.Revision() + 1))
	value.WriteInt64(time.Now().UnixMilli())
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
	err := s.get(&key, &value)
	if err != nil {
		return err
	}
	cid, err := value.ReadInt32()
	if err != nil {
		return err
	}
	if cid != int32(t.ClassId()) {
		return fmt.Errorf("class id not matched %d , %d", cid, t.ClassId())
	}
	rv, err := value.ReadInt64()
	if err != nil {
		return err
	}
	tm, err := value.ReadInt64()
	if err != nil {
		return err
	}
	t.OnTimestamp(tm)
	t.OnRevision(rv)
	t.Read(&value)
	return nil
}

func (s *BadgerLocal) Query(opt core.ListingOpt, stream core.Stream) error {
	if opt.Prefix == nil {
		return fmt.Errorf("query prefix cannot be nil")
	}
	op := badger.DefaultIteratorOptions
	op.Reverse = opt.Reverse
	if op.Reverse && opt.StartCursor == nil {
		opt.StartCursor = append(opt.Prefix, 0xFF)
	}
	if opt.PrefetchValues {
		op.PrefetchValues = true
		op.PrefetchSize = opt.PrefetchSize
	}
	p := opt.Prefix
	seek := p
	if opt.StartCursor != nil {
		seek = opt.StartCursor
	}
	key := BufferProxy{}
	key.NewProxy(BDG_KEY_SIZE)
	value := BufferProxy{}
	value.NewProxy(BDG_VALUE_SIZE)
	limit := -1
	if opt.Limit > 0 {
		limit = opt.Limit
	}
	return s.Db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(op)
		defer it.Close()
		for it.Seek(seek); it.ValidForPrefix(p); it.Next() {
			kv := it.Item()
			key.Clear()
			err := key.Write(kv.Key())
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
			limit--
			if !stream(&key, &value) || limit == 0 {
				break
			}
		}
		return nil
	})
}

func (s *BadgerLocal) Version(key []byte, stream core.Stream) error {
	op := badger.DefaultIteratorOptions
	op.AllVersions = true
	k := BufferProxy{}
	k.NewProxy(BDG_KEY_SIZE)
	v := BufferProxy{}
	v.NewProxy(BDG_VALUE_SIZE)
	return s.Db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(op)
		defer it.Close()
		for it.Seek(key); it.Valid(); it.Next() {
			kv := it.Item()
			if !slices.Equal(kv.Key(), key) {
				break
			}
			k.Clear()
			err := k.Write(kv.Key())
			if err != nil {
				return err
			}
			err = kv.Value(func(val []byte) error {
				v.Clear()
				return v.Write(val)
			})
			if err != nil {
				return err
			}
			k.Flip()
			v.Flip()
			if !stream(&k, &v) {
				break
			}
		}
		return nil
	})
}

func (s *BadgerLocal) set(key *BufferProxy, value *BufferProxy, t core.Persistentable) error {
	if key.Remaining() == 0 || value.Remaining() == 0 {
		return fmt.Errorf("bad key/value")
	}
	return s.Db.Update(func(txn *badger.Txn) error {
		k, _ := key.Read(0)
		v, _ := value.Read(0)
		var riv int64 = 0
		updating := true
		item, err := txn.Get(k)
		if err == nil {
			value.Clear()
			item.Value(func(val []byte) error {
				value.Write(val)
				value.Flip()
				value.ReadInt32()
				v, err := value.ReadInt64()
				if err != nil {
					riv = v
				}
				return nil
			})
		} else { //no record existed
			updating = false
		}
		if riv > 0 && riv != t.Revision() {
			return fmt.Errorf("revison number not matched %d %d", riv, t.Revision())
		}
		if err := txn.Set(k, v); err != nil {
			return err
		}
		if updating {
			return nil
		}
		//update stat total
		se := event.StatEvent{Tag: t.ETag(), Name: event.STAT_TOTAL}
		key.Clear()
		se.WriteKey(key)
		key.Flip()
		value.Clear()
		ck, _ := key.Read(0)
		citem, err := txn.Get(ck)
		if err == nil {
			err = citem.Value(func(val []byte) error {
				return value.Write(val)
			})
			if err != nil {
				return err //rollback
			}
			value.Flip()
			value.ReadInt32()
			value.ReadInt64()
			value.ReadInt64()
			se.Read(value)
			se.Count = se.Count + 1
		} else {
			se.Count = 1
		}
		value.Clear()
		value.WriteInt32(int32(se.ClassId()))
		value.WriteInt64(se.Revision())
		value.WriteInt64(time.Now().UnixMilli())
		se.Write(value)
		value.Flip()
		cv, _ := value.Read(0)
		return txn.Set(ck, cv)

	})
}

func (s *BadgerLocal) get(key *BufferProxy, value *BufferProxy) error {
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

func (s *BadgerLocal) Open() error {
	var opt badger.Options
	if s.InMemory {
		opt = badger.DefaultOptions("").WithInMemory(true)
		if s.LogDisabled {
			opt.Logger = nil
		}
	} else {
		opt = badger.DefaultOptions(s.Path)
		opt.SyncWrites = false
		if s.LogDisabled {
			opt.Logger = nil
		}
	}
	db, err := badger.Open(opt)
	if err != nil {
		return err
	}
	s.Db = db
	return nil
}

func (s *BadgerLocal) Close() error {
	if s.InMemory {
		return s.Db.Close()
	}
	s.Db.Sync()
	return s.Db.Close()
}

func (s *BadgerLocal) GC() error {
	return s.Db.RunValueLogGC(0.7)
}
