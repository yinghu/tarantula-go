package persistence

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gameclustering.com/internal/core"
	badger "github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/ristretto/v2/z"
)

func TestStringKey(t *testing.T) {
	local := BadgerLocal{InMemory: true, LogDisabled: true}
	err := local.Open()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	defer local.Close()
	prefix := "etag:user"
	local.Db.Update(func(txn *badger.Txn) error {
		for i := range 10 {
			key := fmt.Sprintf("%s:%d:uuid%d", prefix, i, i)
			txn.Set([]byte(key), []byte(key))
		}
		return nil
	})
	ct := 0
	local.Db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek([]byte(prefix)); it.ValidForPrefix([]byte(prefix)); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				ct++
				return nil
			})
		}
		return nil
	})
	if ct != 10 {
		t.Errorf("Item should be 10 %d", ct)
	}
	ct = 0
	local.Db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		p := fmt.Sprintf("%s:%d:uuid%d", prefix, 5, 5)
		for it.Seek([]byte(p)); it.ValidForPrefix([]byte(p)); it.Next() {
			ct++
			p = prefix
		}
		return nil
	})
	if ct != 5 {
		t.Errorf("Item should be 5 %d", ct)
	}
	ct = 0
	local.Db.View(func(txn *badger.Txn) error {
		op := badger.IteratorOptions{PrefetchSize: 10, PrefetchValues: false, Reverse: true}
		it := txn.NewIterator(op)

		defer it.Close()
		p := fmt.Sprintf("%s:%d:uuid%d", prefix, 5, 5)
		for it.Seek([]byte(p)); it.ValidForPrefix([]byte(p)); it.Next() {
			ct++
			p = prefix
		}
		return nil
	})
	if ct != 6 {
		t.Errorf("Item should be 6 %d", ct)
	}
}

func TestBufferKey(t *testing.T) {
	local := BadgerLocal{InMemory: true, LogDisabled: true}
	err := local.Open()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	defer local.Close()
	buff := NewBuffer(100)
	local.Db.Update(func(txn *badger.Txn) error {
		for i := range 1000 {
			buff.Clear()
			buff.WriteString("etg")
			buff.WriteString("user")
			buff.WriteInt32(int32(i))
			buff.WriteString(fmt.Sprintf("%s%d", "uuid", i))
			buff.Flip()
			k, _ := buff.Read(0)
			txn.Set(k, k)
		}
		return nil
	})
	ct := 0
	local.Db.View(func(txn *badger.Txn) error {
		for i := range 1000 {
			buff.Clear()
			buff.WriteString("etg")
			buff.WriteString("user")
			buff.WriteInt32(int32(i))
			buff.WriteString(fmt.Sprintf("%s%d", "uuid", i))
			buff.Flip()
			k, _ := buff.Read(0)
			_, err := txn.Get(k)
			if err == nil {
				ct++
			}
		}
		return err
	})
	if ct != 1000 {
		t.Errorf("should be 100 item %d", ct)
	}
	ct = 0
	local.Db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		buff.Clear()
		buff.WriteString("etg")
		buff.WriteString("user")
		buff.Flip()
		prefix, _ := buff.Read(0)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				ct++
				buff.Clear()
				buff.Write(val)
				buff.Flip()
				//s1, _ := buff.ReadString()
				//s2, _ := buff.ReadString()
				//s3, _ := buff.ReadInt32()
				//s4, _ := buff.ReadString()
				//fmt.Printf("%s%s%d%s\n", s1, s2, s3, s4)
				return nil
			})
		}
		return nil
	})
	if ct != 1000 {
		t.Errorf("should be 100 item %d", ct)
	}

	ct = 0
	local.Db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		buff.Clear()
		buff.WriteString("etg")
		buff.WriteString("user")
		buff.Flip()
		prefix, _ := buff.Read(0)
		buff.Clear()
		buff.WriteString("etg")
		buff.WriteString("user")
		buff.WriteInt32(500)
		buff.WriteString(fmt.Sprintf("%s%d", "uuid", 500))
		buff.Flip()
		p, _ := buff.Read(0)
		for it.Seek(p); it.ValidForPrefix(p); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				ct++
				p = prefix
				buff.Clear()
				buff.Write(val)
				buff.Flip()
				//s1, _ := buff.ReadString()
				//s2, _ := buff.ReadString()
				//s3, _ := buff.ReadInt32()
				//s4, _ := buff.ReadString()
				//fmt.Printf("%s%s%d%s\n", s1, s2, s3, s4)
				return nil
			})
		}
		return nil
	})
	if ct != 500 {
		t.Errorf("should be 100 item %d", ct)
	}
	//fmt.Println("Reverse")
	ct = 0
	local.Db.View(func(txn *badger.Txn) error {
		op := badger.IteratorOptions{PrefetchSize: 100, PrefetchValues: false, Reverse: true}
		it := txn.NewIterator(op)
		defer it.Close()
		buff.Clear()
		buff.WriteString("etg")
		buff.WriteString("user")
		buff.Flip()
		prefix, _ := buff.Read(0)
		buff.Clear()
		buff.WriteString("etg")
		buff.WriteString("user")
		buff.WriteInt32(500)
		buff.WriteString(fmt.Sprintf("%s%d", "uuid", 500))
		buff.Flip()
		p, _ := buff.Read(0)
		for it.Seek(p); it.ValidForPrefix(p); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				ct++
				p = prefix
				buff.Clear()
				buff.Write(val)
				buff.Flip()
				//s1, _ := buff.ReadString()
				//s2, _ := buff.ReadString()
				//s3, _ := buff.ReadInt32()
				//s4, _ := buff.ReadString()
				//fmt.Printf("%s%s%d%s\n", s1, s2, s3, s4)
				return nil
			})
		}
		return nil
	})
	if ct != 501 {
		t.Errorf("should be 501 item %d", ct)
	}
}

func TestStreming(t *testing.T) {
	local := BadgerLocal{InMemory: true, LogDisabled: true}
	err := local.Open()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	defer local.Close()
	local.Db.Update(func(txn *badger.Txn) error {
		buff := NewBuffer(100)
		for i := range 1000 {
			buff.Clear()
			buff.WriteString("etg")
			buff.WriteString("user")
			buff.WriteInt32(int32(i))
			buff.WriteString(fmt.Sprintf("%s%d", "uuid", i))
			buff.Flip()
			k, _ := buff.Read(0)
			txn.Set(k, k)
		}
		return nil
	})
	ct := 0
	stream := local.Db.NewStream()
	stream.NumGo = 3
	stream.ChooseKey = func(item *badger.Item) bool {
		buff := NewBuffer(100)
		buff.Clear()
		buff.Write(item.Key())
		buff.Flip()
		//s1, _ := buff.ReadString()
		//s2, _ := buff.ReadString()
		//s3, _ := buff.ReadInt32()
		//s4, _ := buff.ReadString()
		//fmt.Printf("Streaming : %s%s%d%s\n", s1, s2, s3, s4)
		ct++
		return true
	}
	stream.KeyToList = nil
	stream.Send = func(buf *z.Buffer) error {
		return nil
	}
	stream.Orchestrate(context.Background())
	if ct != 1000 {
		t.Errorf("should be 1000 item %d", ct)
	}
}

func TestVersioning(t *testing.T) {
	local := BadgerLocal{InMemory: true, LogDisabled: true}
	err := local.Open()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	defer local.Close()
	k := []byte("key")
	local.Db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(k, []byte("v1"))
		txn.SetEntry(e)

		return nil
	})
	local.Db.Update(func(txn *badger.Txn) error {
		txn.Set(k, []byte("v2"))
		return nil
	})
	local.Db.Update(func(txn *badger.Txn) error {
		txn.Set(k, []byte("v3"))
		return nil
	})
	local.Db.Update(func(txn *badger.Txn) error {
		txn.Set(k, []byte("v4"))
		return nil
	})
	k1 := []byte("key1")
	local.Db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(k1, []byte("v1"))
		txn.SetEntry(e)

		return nil
	})
	local.Db.Update(func(txn *badger.Txn) error {
		txn.Set(k1, []byte("v2"))
		return nil
	})
	local.Db.Update(func(txn *badger.Txn) error {
		txn.Set(k1, []byte("v3"))
		return nil
	})
	local.Db.Update(func(txn *badger.Txn) error {
		txn.Set(k1, []byte("v4"))
		return nil
	})
	var v string
	err = local.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {
			v = string(val)
			return nil
		})
		return nil
	})
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
		return
	}
	if v != "v4" {
		t.Errorf("should be v4 %s", v)
	}
	ct := 0
	local.Db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		opts.AllVersions = true
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(k1); it.Valid(); it.Next() {
			item := it.Item()
			if string(item.Key()) != string(k1) {
				break
			}
			err := item.Value(func(v []byte) error {
				ct++
				//fmt.Printf("key=%s, value=%s, version=%d\n", k1, v, item.Version())
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if ct != 4 {
		t.Errorf("should be 4 item %d", ct)
	}
	local.Version(k, func(k, v core.DataBuffer) bool {
		bk , _ := k.Read(0)
		bv , _ := v.Read(0)
		fmt.Printf("%s %s\n",string(bk),string(bv))
		return true
	})
}

func TestMerging(t *testing.T) {
	local := BadgerLocal{InMemory: true, LogDisabled: true}
	err := local.Open()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	defer local.Close()
	k := []byte("key")
	local.Db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(k, []byte("v1"))
		txn.SetEntry(e)

		return nil
	})
	m := local.Db.GetMergeOperator(k, func(old, new []byte) []byte {
		old = append(old, new...)
		return old
	}, 100*time.Millisecond)
	defer m.Stop()
	m.Add([]byte("v2"))
	m.Add([]byte("v3"))
	m.Add([]byte("v4"))
	res, _ := m.Get()
	v := string(res)
	if v != "v1v2v3v4" {
		t.Errorf("should be v1v2v3v4 %s", v)
	}
	//fmt.Printf("Merged %s\n", v)
}
