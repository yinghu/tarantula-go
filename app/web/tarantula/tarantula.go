package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/tarantula/presence"
)

func main() {
	core.CreateTestLog()
	m := MemberlistManager{}
	m.Seed = []string{"192.168.1.11", "192.168.1.6", "192.168.1.7"}
	err := m.Start()
	if err != nil {
		core.AppLog.Printf("no cluster can join %s", err.Error())
		return
	}
	go presence.Start(&m)
	//m.Cs = &m
	//go m.DataServiceProvider.Start()
	go m.ShutdownHook()
	select {}
}
