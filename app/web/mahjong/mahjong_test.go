package main

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
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
	dealer := (mt.Pts - 1) % 4
	dz := len(mt.Players[dealer].Hand.Tiles)
	if dz != 14 {
		t.Errorf("dealer hand should be 14 %d", dz)
	}
	//fmt.Printf("F Hand %v\n",mt.Players[dealer].Tiles)
	mt.Claim(dealer)
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
	dealer := (mt.Pts - 1) % 4
	dz := len(mt.Players[dealer].Hand.Tiles)
	if dz != 14 {
		t.Errorf("dealer hand should be 14 %d", dz)
	}
	mt.Claim(dealer)
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
			fmt.Printf("X Tiles %v\n", mt.Players[i].Tiles)
			err = mt.Draw(i)
			if err != nil {
				t.Errorf("shoud not be error %s", err.Error())
			}
			fmt.Printf("Y Tiles %v\n", mt.Players[i].Tiles)
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
	me := MahjongEvent{Token: MahjongPlayToken{Cmd: CMD_DICE, SystemId: 100}, SystemId: 101}
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

func TestMahjongTableIndex(t *testing.T) {
	ix := make(map[int64]*MahjongTable)
	mt := MahjongTable{}
	mt.New()
	ix[10] = &mt
	mx := ix[10]
	mx.Dice()
	mx.Deal()
	tx := mt.Players[SEAT_E].Hand.Tiles[0]
	fmt.Printf("%v\n", mt.Players[SEAT_E].Hand.Tiles)
	fmt.Printf("%v\n", mx.Players[SEAT_E].Hand.Tiles)
	mx.Discharge(SEAT_E, tx.Seq)
	fmt.Printf("%v\n", mt.Players[SEAT_E].Hand.Tiles)
	fmt.Printf("%v\n", mx.Players[SEAT_E].Hand.Tiles)

}
