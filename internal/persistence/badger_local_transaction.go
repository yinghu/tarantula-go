package persistence

import badger "github.com/dgraph-io/badger/v4"

type BadgerLocalTransaction struct {
	txn *badger.Txn
}

func (t *BadgerLocalTransaction) Get(key []byte) ([]byte, error) {
	item, err := t.txn.Get(key)
	if err != nil {
		return nil, err
	}
	v := make([]byte, 0)
	err = item.Value(func(val []byte) error {
		v = append(v, val...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (t *BadgerLocalTransaction) Set(key []byte, value []byte) error {
	return t.txn.Set(key,value)
}

func (t *BadgerLocalTransaction) Del(key []byte) error {
	return t.txn.Delete(key)
}

func (t *BadgerLocalTransaction) Commit() error {
	return t.txn.Commit()
}

func (t *BadgerLocalTransaction) Rollback() {
	t.txn.Discard()
}
