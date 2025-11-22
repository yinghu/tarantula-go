package main

import (
	"time"

	"gameclustering.com/internal/event"
)

type MahjongTimeout interface {
	Start()
	Stop()
	OId() int64
}

type MahjongTimeoutObj struct {
	event.EventObj
	T time.Timer
}
