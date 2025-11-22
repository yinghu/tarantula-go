package main

import (
	"gameclustering.com/internal/event"
)

type MahjongTimeout interface {
	Start()
	Stop()
	OId() int64
}

type MahjongTimeoutObj struct {
	event.EventObj
	Commited bool
	//T        time.Timer
}
