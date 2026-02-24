package main

import (
	"os"
	"os/signal"
	"syscall"

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
	//shutdown hook
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	core.AppLog.Warn().Msg("shutdown signal to be triggered")
	m.ShutdownHook()
	signal.Stop(sigs)
	close(sigs)
	os.Exit(0)
}
