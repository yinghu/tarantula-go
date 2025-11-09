package main

import (
	"gameclustering.com/internal/mj"
)

type MahjongPlayer struct {
	SystemId int64
	Seat     string
	mj.Hand
	Auto bool
	B []mj.Tile //bamboo
	C []mj.Tile //character
	D []mj.Tile //dots
	HE []mj.Tile //east
	HS []mj.Tile //south
	HW []mj.Tile //west
	HN []mj.Tile //north
	R []mj.Tile //red
	G []mj.Tile //green
	W []mj.Tile //white
}

func NewPlayer(seat string) MahjongPlayer {
	mp := MahjongPlayer{Seat: seat, Auto: true}
	mp.Hand = mj.Hand{}
	mp.Hand.New()
	mp.B = make([]mj.Tile, 0)
	mp.C = make([]mj.Tile, 0)
	mp.D = make([]mj.Tile, 0)
	mp.HE = make([]mj.Tile, 0)
	mp.HS = make([]mj.Tile, 0)
	mp.HW = make([]mj.Tile, 0)
	mp.HN = make([]mj.Tile, 0)
	mp.R = make([]mj.Tile, 0)
	mp.G = make([]mj.Tile, 0)
	mp.W = make([]mj.Tile, 0)
	return mp
}
