package main

import (
	"testing"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/tarantula/protocol"
	badger "github.com/dgraph-io/badger/v4"
	hash "github.com/spaolacci/murmur3"
)

func TestKeyIndex2(t *testing.T) {
	core.CreateTestLog()
	local := persistence.BadgerLocal{InMemory: true, LogDisabled: true}
	err := local.Open()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	defer local.Close()
	dsp := DataServiceProvider{Local: &local}
	key := []byte("key1")
	kw := KeyIndex{Prefix: hash.Sum32(key), Key: key, Header: &protocol.Header{Revision: 1, FactoryId: 1, ClassId: 1, Timestamp: 1, Size: 100}}
	err = dsp.saveKeyIndex(&kw)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	kr := KeyIndex{Prefix: kw.Prefix, Key: kw.Key, Header: &protocol.Header{FactoryId: 1, ClassId: 1}}
	err = dsp.loadKeyIndex(&kr)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	if kw.Header.Revision != kr.Header.Revision {
		t.Errorf("should be same %d %d", kw.Header.Revision, kr.Header.Revision)
	}
	if kw.Header.FactoryId != kr.Header.FactoryId {
		t.Errorf("should be same %d %d", kw.Header.FactoryId, kr.Header.FactoryId)
	}
	if kw.Header.ClassId != kr.Header.ClassId {
		t.Errorf("should be same %d %d", kw.Header.ClassId, kr.Header.ClassId)
	}
	if kw.Header.Timestamp != kr.Header.Timestamp {
		t.Errorf("should be same %d %d", kw.Header.Timestamp, kr.Header.Timestamp)
	}
	if kw.Header.Size != kr.Header.Size {
		t.Errorf("should be same %d %d", kw.Header.Size, kr.Header.Size)
	}
}
func TestKeyIndex3(t *testing.T) {
	core.CreateTestLog()
	local := persistence.BadgerLocal{InMemory: true, LogDisabled: true}
	err := local.Open()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	defer local.Close()
	in := protocol.Data{Key: []byte("key"), Value: []byte("value123"), Header: &protocol.Header{FactoryId: 1, ClassId: 1, Revision: 1, Timestamp: time.Now().UnixMilli(), Size: 8}}
	setData := SetData{Data: &in}
	ki := setData.KeyIndex()
	core.AppLog.Debug().Msgf("Key index : %v", ki)
	k, v, err := ki.Kv()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	ak, err := setData.K()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	err = local.Db.Update(func(txn *badger.Txn) error {
		err := txn.Set(k, v)
		if err != nil {
			return err
		}
		return txn.Set(ak, setData.Value)
	})
	if err != nil {
		t.Errorf("should be commited %s", err.Error())
	}
	kx := KeyIndex{Header: &protocol.Header{}}
	err = local.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			return kx.V(val)

		})
		if err != nil {
			return err
		}
		item, err = txn.Get(ak)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			core.AppLog.Debug().Msgf("value loaded : %s", string(val))
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Errorf("should be loaded %s", err.Error())
	}
	core.AppLog.Debug().Msgf("loaded %v", kx)

}
