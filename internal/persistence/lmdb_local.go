package persistence

import (
	"fmt"
	"runtime"

	"github.com/PowerDNS/lmdb-go/lmdb"
	//"github.com/PowerDNS/lmdb-go/lmdbscan"
)

const (
	MB_BYTES      int64 = 1048576
	READER_NUMBER int   = 100
	MAX_DB_NUMBER int   = 1024
)

type LMDBLocal struct {
	env       *lmdb.Env
	Path      string
	MapSizeMb int64
	Readers   int
}

func (db *LMDBLocal) Open() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	env, err := lmdb.NewEnv()
	if err != nil {
		fmt.Printf("lmd error %s\n", err.Error())
		panic(err)
	}
	db.env = env
	env.SetMapSize(db.MapSizeMb * MB_BYTES)
	env.SetMaxReaders(db.Readers)
	env.SetMaxDBs(MAX_DB_NUMBER)
	err = db.env.Open(db.Path, 0, 0644)
	if err != nil {
		fmt.Printf("lmd error %s\n", err.Error())
		panic(err)
	}
}

func (db *LMDBLocal) Put(dbName string, key, value []byte) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return db.env.Update(func(txn *lmdb.Txn) error {
		dbi, err := txn.CreateDBI(dbName)
		if err != nil {
			return err
		}
		return txn.Put(dbi, key, value, 0)
	})
}

func (db *LMDBLocal) PutEdge(dbName string, label string, key, value []byte) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return db.env.Update(func(txn *lmdb.Txn) error {
		dbi, err := txn.CreateDBI(fmt.Sprintf("%s#%s", dbName, label))
		if err != nil {
			return err
		}
		return txn.Put(dbi, key, value, lmdb.DupSort)
	})
}

func (db *LMDBLocal) Get(dbName string, key []byte) ([]byte, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	value := make([]byte, 0)
	err := db.env.View(func(txn *lmdb.Txn) error {
		dbi, err := txn.CreateDBI(dbName)
		if err != nil {
			return err
		}
		v, err := txn.Get(dbi, key)
		if err != nil {
			return err
		}
		value = append(value, v...)
		return nil
	})
	return value, err
}

func (db *LMDBLocal) Del(dbName string, key []byte) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return db.env.Update(func(txn *lmdb.Txn) error {
		dbi, err := txn.CreateDBI(dbName)
		if err != nil {
			return err
		}
		return txn.Del(dbi, key, key)
	})
}

func (db *LMDBLocal) Close() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	db.env.Sync(true)
	db.env.Close()
}
