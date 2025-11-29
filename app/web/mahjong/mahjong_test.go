package main

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/mj"
)

type SampleCallback struct {
}

func (s *SampleCallback) OnError(e event.Event, err error) {
	fmt.Printf("On event error %v %s\n", e, err.Error())
}

func (s *SampleCallback) OnEvent(e event.Event) {
	fmt.Printf("On event %v\n", e)
}

func (s *SampleCallback) Push(e event.Event) {

}

func TestMahjongTable(t *testing.T) {
	core.CreateTestLog()
	mt := MahjongTable{Pusher: &SampleCallback{}}
	mt.New()
	mt.Sit(1, SEAT_E)
	mt.Sit(2, SEAT_S)
	mt.Sit(3, SEAT_W)
	mt.Sit(4, SEAT_N)
	mt.Dice()
	mt.Deal()
	dealer := (mt.Pts() - 1) % 4
	dz := len(mt.Players[dealer].Hand.Tiles)
	if dz != 14 {
		t.Errorf("dealer hand should be 14 %d", dz)
	}
	//fmt.Printf("F Hand %v\n",mt.Players[dealer].Tiles)
	//mt.Claim(dealer)
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
	mt.New()
	//mt.Sit(1, SEAT_E)
	//mt.Sit(2, SEAT_S)
	//mt.Sit(3, SEAT_W)
	//mt.Sit(4, SEAT_N)
	mt.Dice()
	mt.Deal()
	dealer := (mt.Pts() - 1) % 4
	dz := len(mt.Players[dealer].Hand.Tiles)
	if dz != 14 {
		t.Errorf("dealer hand should be 14 %d", dz)
	}
	//mt.Claim(dealer)
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
			//fmt.Printf("X Tiles %v\n", mt.Players[i].Tiles)
			err = mt.Draw(i)
			if err != nil {
				t.Errorf("shoud not be error %s", err.Error())
			}
			//fmt.Printf("Y Tiles %v\n", mt.Players[i].Tiles)
			hz := len(mt.Players[i].Tiles)
			if hz != 14 {
				t.Errorf("hand size should be 14 %d", hz)
			}
		}
	}
	tx := mt.Players[dealer].Hand.Tiles[0]
	err = mt.Discharge(dealer, tx.Seq)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	sz := len(mt.Players[dealer].Hand.Tiles)
	if sz != 13 {
		t.Errorf("hand size should be 13 %d", sz)
	}
	err = mt.Draw(dealer)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	sz = len(mt.Players[dealer].Hand.Tiles)
	if sz != 14 {
		t.Errorf("hand size should be 14 %d", sz)
	}
}

func TestMahjongToken(t *testing.T) {
	lis := SampleCallback{}
	me := NewMahjongEvent(101, MahjongPlayToken{Cmd: CMD_DICE, SystemId: 100})
	me.Callback = &lis
	buff := core.NewBuffer(100)
	me.Outbound(buff)
	mx := MahjongEvent{Token: MahjongPlayToken{}}
	mx.Callback = &lis
	buff.Flip()
	mx.Inbound(buff)
	if me.Token != mx.Token {
		t.Errorf("token should be same %v %v", me.Token, mx.Token)
	}
}

// func cmp(a, b mj.Tile) int {
// return a.Seq - b.Seq
// }
func TestClassic1(t *testing.T) {

	b1 := mj.FromS(mj.BAMBOO1)
	b2 := mj.FromS(mj.BAMBOO1)

	b3 := mj.FromS(mj.BAMBOO2)

	b4 := mj.FromS(mj.BAMBOO3)
	b5 := mj.FromS(mj.BAMBOO4)
	b6 := mj.FromS(mj.BAMBOO6)

	b7 := mj.FromS(mj.BAMBOO6)
	b8 := mj.FromS(mj.BAMBOO7)
	b9 := mj.FromS(mj.BAMBOO7)

	c7 := mj.FromS(mj.BAMBOO8)
	c8 := mj.FromS(mj.BAMBOO8)

	c9 := mj.FromS(mj.DOTS5)

	p1 := mj.FromS(mj.DOTS6)
	p2 := mj.FromS(mj.DOTS7)

	tiles := []mj.Tile{b1, b2, b3, b4, b5, b6, b7, b8, b9, c7, c8, c9, p1, p2}
	//checkList(tiles)
	if len(tiles) != 14 {
		t.Errorf("hand size should be 14 %d", len(tiles))
	}
	h := mj.Hand{}
	h.New()
	h.Tiles = append(h.Tiles, tiles...)
	//slices.SortFunc(h.Tiles, cmp)
	cm := mj.ClassicMahjong{}
	cm.New()
	claimed := cm.Mahjong(&h)
	if !claimed {
		t.Errorf("should be claimed %v", claimed)
	}
	//B1,B2,B3,B2,B3,B4,B3,B4,B5,B4,B5,B6,B1,B1
}

func TestClassic2(t *testing.T) {

	b1 := mj.FromS(mj.BAMBOO5)
	b2 := mj.FromS(mj.BAMBOO6)
	b3 := mj.FromS(mj.BAMBOO7)

	b4 := mj.FromS(mj.BAMBOO7)
	b5 := mj.FromS(mj.BAMBOO8)
	b6 := mj.FromS(mj.BAMBOO9)

	b7 := mj.FromS(mj.CHARACTER3)
	b8 := mj.FromS(mj.CHARACTER4)
	b9 := mj.FromS(mj.CHARACTER5)

	c7 := mj.FromS(mj.CHARACTER6)
	c8 := mj.FromS(mj.CHARACTER6)

	c9 := mj.FromS(mj.RED)

	p1 := mj.FromS(mj.RED)
	p2 := mj.FromS(mj.RED)

	tiles := []mj.Tile{b1, b2, b3, b4, b5, b6, b7, b8, b9, c7, c8, c9, p1, p2}
	h := mj.Hand{}
	h.New()
	h.Tiles = append(h.Tiles, tiles...)
	cm := mj.ClassicMahjong{}
	cm.New()
	claimed := cm.Mahjong(&h)
	if !claimed {
		t.Errorf("should be claimed %v", claimed)
	}
	//B1,B2,B3,B2,B3,B4,B3,B4,B5,B4,B5,B6,B1,B1
}
