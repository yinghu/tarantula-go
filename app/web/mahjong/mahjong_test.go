package main

import (
	"testing"
)

func TestMahjongTable(t *testing.T) {
	mt := MahjongTable{}
	mt.New()
	mt.Sit(1, SEAT_E)
	mt.Sit(2, SEAT_S)
	mt.Sit(3, SEAT_W)
	mt.Sit(4, SEAT_N)
	mt.Dice()
	mt.Deal()
	dealer := (mt.Pts - 1) % 4
	dz := len(mt.Players[dealer].Hand.Tiles)
	if dz != 14 {
		t.Errorf("dealer hand should be 14 %d", dz)
	}
	for i := range 4 {
		if i != dealer {
			pz := len(mt.Players[i].Tiles)
			if pz != 13 {
				t.Errorf("player hand should be 12 %d", pz)
			}
		}
	}
}
