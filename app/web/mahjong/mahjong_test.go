package main

import (
	"fmt"
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
	for i := range 4 {
		fmt.Printf("HND SIZE : %d\n", len(mt.Players[i].Hand.Tiles)+ len(mt.Players[i].Flowers))
	}
}
