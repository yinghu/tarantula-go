package persistence

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/core"
	//"github.com/PowerDNS/lmdb-go/lmdb"
)

func TestSimple(t *testing.T) {
	lm := LMDBLocal{Path: "/home/yinghu/lmdb", MapSizeMb: 100, Readers: 100}
	lm.Open()
	defer lm.Close()

	err := lm.Put("test", []byte("key2"), []byte("value"))
	if err != nil {
		fmt.Printf("put error %s\n", err.Error())
		return
	}
	v, _ := lm.Get("test", []byte("key"))
	fmt.Printf("value : %s\n", string(v))
}

func TestAuto(t *testing.T) {
	core.CreateTestLog()
	lm := LMDBLocalAuto{Path: "/home/yinghu/auto", MapSizeMb: 100, Readers: 100, MaxDbs: 1024, BatchSzie: 16}
	lm.Open()
	defer lm.Env.TerminateSync()

	//err := lm.Put("test", []byte("key2"), []byte("value"))
	//if err != nil {
	//fmt.Printf("put error %s\n", err.Error())
	//return
	//}
	//v, _ := lm.Get("test", []byte("key"))
	//fmt.Printf("value : %s\n", string(v))
}
