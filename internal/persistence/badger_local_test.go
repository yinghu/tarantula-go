package persistence

import (
	"fmt"
	"testing"

	badger "github.com/dgraph-io/badger/v4"
)

func TestLocalData(t *testing.T) {
	local := BadgerLocal{InMemory: true, Path: "/home/yinghu/local/t1"}
	err := local.Open()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	defer local.Close()
	prefix := "etag:user"
	local.Db.Update(func(txn *badger.Txn) error {
		for i := range 10 {
			key := fmt.Sprintf("%s:%d:uuid%d", prefix, i, i)
			fmt.Printf("key : %s\n", key)
			txn.Set([]byte(key), []byte(key))
		}
		return nil
	})
	local.Db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek([]byte(prefix)); it.ValidForPrefix([]byte(prefix)); it.Next() {
			item := it.Item()
			fmt.Printf("it : %s\n", string(item.Key()))
		}
		return nil
	})
	local.Db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		p := fmt.Sprintf("%s:%d:uuid%d", prefix, 5, 5)
		fmt.Printf("seek : %s\n", p)
		for it.Seek([]byte(p)); it.ValidForPrefix([]byte(p)); it.Next() {
			item := it.Item()
			fmt.Printf("range : %s\n", string(item.Key()))
			p = prefix
		}
		return nil
	})

	local.Db.View(func(txn *badger.Txn) error {
		op := badger.IteratorOptions{PrefetchSize: 10, PrefetchValues: false, Reverse: true}
		it := txn.NewIterator(op)

		defer it.Close()
		p := fmt.Sprintf("%s:%d:uuid%d", prefix, 5, 5)
		fmt.Printf("reverse : %s\n", p)
		for it.Seek([]byte(p)); it.ValidForPrefix([]byte(p)); it.Next() {
			item := it.Item()
			fmt.Printf("reverse : %s\n", string(item.Key()))
			p = prefix
		}
		return nil
	})
}
