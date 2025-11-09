package main

import (
	"testing"
)

func TestMahjongTable(t *testing.T) {
	mt := MahjongTable{}
	mt.Reset()
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
	//fmt.Printf("F Hand %v\n",mt.Players[dealer].Tiles)
	mt.Mahjong(dealer)
	//fmt.Printf("X Hand %v\n",mt.Players[dealer].Tiles)
	err := mt.Draw(dealer)
	if err == nil {
		t.Errorf("should be error")
	}
	for i := range 4 {
		if i != dealer {
			pz := len(mt.Players[i].Tiles)
			if pz != 13 {
				t.Errorf("player hand should be 13 %d", pz)
			}
			err = mt.Draw(i)
			if err != nil {
				t.Errorf("shoud not be error %s", err.Error())
			}
			hz := len(mt.Players[i].Tiles)
			if hz != 14 {
				t.Errorf("hand size should be 14 %d", hz)
			}
		}
	}
}

func TestMahjongAutoTable(t *testing.T) {
	mt := MahjongTable{}
	mt.Reset()
	//mt.Sit(1, SEAT_E)
	//mt.Sit(2, SEAT_S)
	//mt.Sit(3, SEAT_W)
	//mt.Sit(4, SEAT_N)
	mt.Dice()
	mt.Deal()
	dealer := (mt.Pts - 1) % 4
	dz := len(mt.Players[dealer].Hand.Tiles)
	if dz != 14 {
		t.Errorf("dealer hand should be 14 %d", dz)
	}
	//fmt.Printf("F Hand %v\n",mt.Players[dealer].Tiles)
	mt.Mahjong(dealer)
	//fmt.Printf("X Hand %v\n",mt.Players[dealer].Tiles)
	err := mt.Draw(dealer)
	if err == nil {
		t.Errorf("should be error")
	}
	for i := range 4 {
		if i != dealer {
			pz := len(mt.Players[i].Tiles)
			if pz != 13 {
				t.Errorf("player hand should be 13 %d", pz)
			}
			err = mt.Draw(i)
			if err != nil {
				t.Errorf("shoud not be error %s", err.Error())
			}
			hz := len(mt.Players[i].Tiles)
			if hz != 14 {
				t.Errorf("hand size should be 14 %d", hz)
			}
		}
	}
}
