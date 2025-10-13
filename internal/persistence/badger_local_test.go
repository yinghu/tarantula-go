package persistence

import (
	"fmt"
	"testing"

	badger "github.com/dgraph-io/badger/v4"
)

func TestStringKey(t *testing.T) {
	local := BadgerLocal{InMemory: true}
	err := local.Open()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	defer local.Close()
	prefix := "etag:user"
	local.Db.Update(func(txn *badger.Txn) error {
		for i := range 10 {
			key := fmt.Sprintf("%s:%d:uuid%d", prefix, i, i)
			//fmt.Printf("key : %s\n", key)
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
		//fmt.Printf("seek : %s\n", p)
		for it.Seek([]byte(p)); it.ValidForPrefix([]byte(p)); it.Next() {
			//item := it.Item()
			//fmt.Printf("range : %s\n", string(item.Key()))
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
		//fmt.Printf("reverse : %s\n", p)
		for it.Seek([]byte(p)); it.ValidForPrefix([]byte(p)); it.Next() {
			//item := it.Item()
			//fmt.Printf("reverse : %s\n", string(item.Key()))
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
	local := BadgerLocal{InMemory: true}
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
				s1, _ := buff.ReadString()
				s2, _ := buff.ReadString()
				s3, _ := buff.ReadInt32()
				s4, _ := buff.ReadString()
				fmt.Printf("%s%s%d%s\n", s1, s2, s3, s4)
				return nil
			})
		}
		return nil
	})
	if ct != 1000 {
		t.Errorf("should be 100 item %d", ct)
	}
}
