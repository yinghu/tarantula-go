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
	if me.Token != mx.Token{
		t.Errorf("token should be same %v %v", me.Token,mx.Token)
	}
}
