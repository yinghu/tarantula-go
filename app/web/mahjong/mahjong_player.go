package main

import "gameclustering.com/internal/mj"

type MahjongPlayer struct {
	SystemId int64
	Seat     string
	mj.Hand
}

func NewPlayer() MahjongPlayer {
	mp := MahjongPlayer{}
	mp.Hand = mj.Hand{}
	mp.Hand.New()
	return mp
}
