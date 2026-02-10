package main

import (
	"fmt"
	"os"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/tarantula/data"
	"gameclustering.com/tarantula/presence"
)

func main() {
	core.CreateTestLog()
	m := MemberlistManager{}
	m.Seed = []string{"192.168.1.11", "192.168.1.6", "192.168.1.7"}
	err := m.Start()
	if err != nil {
		core.AppLog.Printf("no cluster can join %s\n", err.Error())
		return
	}
	go presence.Start(&m)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	path := fmt.Sprintf("%s/%s", homeDir, "tarantula")
	core.AppLog.Printf("check path %s\n",path)
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return
	}
	lmdb := persistence.LMDBLocal{Path: path, MapSizeMb: 10, Readers: 100}
	dsp := data.DataServiceProvider{Db: &lmdb,Cs: &m}
	go dsp.Start()
	go m.ShutdownHook()
	select {}
}
