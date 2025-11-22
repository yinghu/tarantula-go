package main

type MahjongTimeout interface {
	Start()
	Stop()
	OId() int64
}
