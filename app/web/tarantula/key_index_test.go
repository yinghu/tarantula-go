package main

import (
	"testing"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"

	badger "github.com/dgraph-io/badger/v4"
	hash "github.com/spaolacci/murmur3"
)

func TestKeyIndex(t *testing.T) {
	core.CreateTestLog()
	slaves := make([]core.Node, 0, 3)
	slaves = append(slaves, core.Node{RingToken: hash.Sum32([]byte("slave1"))})
	slaves = append(slaves, core.Node{RingToken: hash.Sum32([]byte("slave2"))})
	slaves = append(slaves, core.Node{RingToken: hash.Sum32([]byte("slave3"))})
	k := KeyIndex{Prefix: hash.Sum32([]byte("key")), Master: core.Node{RingToken: hash.Sum32([]byte("master"))}, Slaves: slaves, Key: []byte("hello")}
	kBuff := core.NewBuffer(100)
	vBuff := core.NewBuffer(100)
	k.WriteKey(kBuff)
	k.Write(vBuff)
	kBuff.Flip()
	vBuff.Flip()
	var ak KeyIndex
	ak.ReadKey(kBuff)
	ak.Read(vBuff)
	if ak.Prefix != k.Prefix {
		t.Errorf("should be same %d %d", ak.Prefix, k.Prefix)
	}
	//core.AppLog.Debug().Msgf("prefix %d %d", ak.Prefix, k.Prefix)
	if ak.Master.RingToken != k.Master.RingToken {
		t.Errorf("should be same %d %d", ak.Master.RingToken, k.Master.RingToken)
	}
	//core.AppLog.Debug().Msgf("master %d %d", ak.Master.RingToken, k.Master.RingToken)
	if string(ak.Key) != string(k.Key) {
		t.Errorf("should be same %s %s", string(ak.Key), string(k.Key))
	}
	//core.AppLog.Debug().Msgf("key %s %s", string(ak.Key), string(k.Key))
	for i, n := range k.Slaves {
		if n.RingToken != ak.Slaves[i].RingToken {
			t.Errorf("should be same %d %d", n.RingToken, ak.Slaves[i].RingToken)
		}
		//core.AppLog.Debug().Msgf("slave %d %d %d", i, n.RingToken, ak.Slaves[i].RingToken)
	}
}

func TestKeyIndex2(t *testing.T) {
	local := persistence.BadgerLocal{InMemory: true, LogDisabled: true}
	err := local.Open()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	defer local.Close()
	core.CreateTestLog()
	slaves := make([]core.Node, 0, 3)
	slaves = append(slaves, core.Node{RingToken: hash.Sum32([]byte("slave1"))})
	slaves = append(slaves, core.Node{RingToken: hash.Sum32([]byte("slave2"))})
	slaves = append(slaves, core.Node{RingToken: hash.Sum32([]byte("slave3"))})
	k := KeyIndex{Prefix: hash.Sum32([]byte("key")), Master: core.Node{RingToken: hash.Sum32([]byte("master"))}, Slaves: slaves, Key: []byte("hello")}
	kBuff := core.NewBuffer(100)
	vBuff := core.NewBuffer(100)
	k.WriteKey(kBuff)
	k.Write(vBuff)
	kBuff.Flip()
	vBuff.Flip()
	key, _ := kBuff.Read(0)
	if err := local.Db.Update(func(txn *badger.Txn) error {
		value, _ := vBuff.Read(0)
		return txn.Set(key, value)
	}); err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	//read from db
	kBuff.Clear()
	vBuff.Clear()
	if err := local.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			if err := vBuff.Write(val); err != nil {
				return err
			}
			return nil
		})
	}); err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	kBuff.Write(key)
	kBuff.Flip()
	vBuff.Flip()
	ak := KeyIndex{}
	ak.ReadKey(kBuff)
	ak.Read(vBuff)
	if ak.Prefix != k.Prefix {
		t.Errorf("should be same %d %d", ak.Prefix, k.Prefix)
	}
	//core.AppLog.Debug().Msgf("prefix %d %d", ak.Prefix, k.Prefix)
	if ak.Master.RingToken != k.Master.RingToken {
		t.Errorf("should be same %d %d", ak.Master.RingToken, k.Master.RingToken)
	}
	//core.AppLog.Debug().Msgf("master %d %d", ak.Master.RingToken, k.Master.RingToken)
	if string(ak.Key) != string(k.Key) {
		t.Errorf("should be same %s %s", string(ak.Key), string(k.Key))
	}
	//core.AppLog.Debug().Msgf("key %s %s", string(ak.Key), string(k.Key))
	for i, n := range k.Slaves {
		if n.RingToken != ak.Slaves[i].RingToken {
			t.Errorf("should be same %d %d", n.RingToken, ak.Slaves[i].RingToken)
		}
		//core.AppLog.Debug().Msgf("slave %d %d %d", i, n.RingToken, ak.Slaves[i].RingToken)
	}
}

func TestKeyIndex3(t *testing.T) {
	local := persistence.BadgerLocal{InMemory: true, LogDisabled: true}
	err := local.Open()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	defer local.Close()
	core.CreateTestLog()
	dsp := DataServiceProvider{Local: &local}
	slaves := make([]core.Node, 0, 3)
	slaves = append(slaves, core.Node{RingToken: hash.Sum32([]byte("slave1"))})
	slaves = append(slaves, core.Node{RingToken: hash.Sum32([]byte("slave2"))})
	slaves = append(slaves, core.Node{RingToken: hash.Sum32([]byte("slave3"))})
	k := KeyIndex{Prefix: hash.Sum32([]byte("key")), Master: core.Node{RingToken: hash.Sum32([]byte("master"))}, Slaves: slaves, Key: []byte("hello")}
	err = dsp.SaveKeyIndex(&k)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	ak := KeyIndex{Prefix: hash.Sum32([]byte("key")),Key: []byte("hello")}
	err = dsp.LoadKeyIndex(&ak)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	if ak.Prefix != k.Prefix {
		t.Errorf("should be same %d %d", ak.Prefix, k.Prefix)
	}
	//core.AppLog.Debug().Msgf("prefix %d %d", ak.Prefix, k.Prefix)
	if ak.Master.RingToken != k.Master.RingToken {
		t.Errorf("should be same %d %d", ak.Master.RingToken, k.Master.RingToken)
	}
	//core.AppLog.Debug().Msgf("master %d %d", ak.Master.RingToken, k.Master.RingToken)
	if string(ak.Key) != string(k.Key) {
		t.Errorf("should be same %s %s", string(ak.Key), string(k.Key))
	}
	//core.AppLog.Debug().Msgf("key %s %s", string(ak.Key), string(k.Key))
	for i, n := range k.Slaves {
		if n.RingToken != ak.Slaves[i].RingToken {
			t.Errorf("should be same %d %d", n.RingToken, ak.Slaves[i].RingToken)
		}
		//core.AppLog.Debug().Msgf("slave %d %d %d", i, n.RingToken, ak.Slaves[i].RingToken)
	}
}
