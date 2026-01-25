package persistence

import (
	"fmt"
	"runtime"

	"github.com/PowerDNS/lmdb-go/lmdb"
	//"github.com/PowerDNS/lmdb-go/lmdbscan"
)

type LMDBLocal struct {
	e *lmdb.Env
}

func (db *LMDBLocal) Open() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	env , err := lmdb.NewEnv()
	if err != nil {
		fmt.Printf("lmd error %s\n",err.Error())
		panic(err)
	}
	db.e = env
	err = db.e.Open("/home/yinghu/lmdb",0,0644)
	if err != nil {
		fmt.Printf("lmd error %s\n",err.Error())
		panic(err)
	}
}

func (db *LMDBLocal) Close() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	db.e.Sync(true)
	db.e.Close()
}
