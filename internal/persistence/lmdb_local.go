package persistence

import (
	"runtime"

	"github.com/PowerDNS/lmdb-go/lmdb"
	//"github.com/PowerDNS/lmdb-go/lmdbscan"
)

type LMDBLocal struct {
}

func (db *LMDBLocal) Open() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	_, err := lmdb.NewEnv()
	if err != nil {
		panic(err)
	}

}
