package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/tarantula/presence"
)

func main() {
	core.CreateTestLog()
	presence.Start()
	m := MemberlistManager{}
	m.Seed = []string{"192.168.1.11", "192.168.1.6"}
	err := m.Start()
	if err != nil {
		core.AppLog.Printf("no cluster can join %s\n", err.Error())
		return
	}
	select {}
}
