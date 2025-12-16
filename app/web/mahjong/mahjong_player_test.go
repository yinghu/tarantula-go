package main

import (
	"testing"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/mj"
)

type SamplePusher struct {
}

func (p *SamplePusher) Push(e event.Event) {

}

func TestPlayerSetup(t *testing.T) {
	core.CreateTestLog()
	mp := NewPlayer(1, true, &SampleCallback{})
	if !mp.Auto {
		t.Errorf("default should be auto %v", mp.Auto)
	}
	if mp.FlowerExcluded {
		t.Errorf("flower should be included %v", mp.FlowerExcluded)
	}
	if mp.Seat != 1 {
		t.Errorf("seat should be 1 %d", mp.Seat)
	}
	if !mp.Sorting {
		t.Errorf("hand should be shorting %v", mp.Sorting)
	}
	mp.AppendForTest(mj.FromS(mj.BAMBOO1), false)
	mp.AppendForTest(mj.FromS(mj.BAMBOO2), false)
	mp.AppendForTest(mj.FromS(mj.BAMBOO3), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER1), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER1), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER1), false)
	mp.AppendForTest(mj.FromS(mj.DOTS1), false)
	mp.AppendForTest(mj.FromS(mj.DOTS1), false)
	mp.AppendForTest(mj.FromS(mj.DOTS1), false)
	mp.AppendForTest(mj.FromS(mj.DOTS1), false)
	mp.AppendForTest(mj.FromS(mj.RED), false)
	mp.AppendForTest(mj.FromS(mj.RED), false)
	mp.AppendForTest(mj.FromS(mj.F_AUTUMN), false)
	if mp.TileSize() != 13 {
		t.Errorf("hand size should be 13 %d", mp.TileSize())
	}
	if len(mp.B) != 3 {
		t.Errorf("bamboo seg size should be 3 %d", len(mp.B))
	}
	if len(mp.C) != 3 {
		t.Errorf("characters seg size should be 3 %d", len(mp.C))
	}
	if len(mp.D) != 4 {
		t.Errorf("DOTS seg size should be 4 %d", len(mp.D))
	}
	if len(mp.R) != 2 {
		t.Errorf("red seg size should be 2 %d", len(mp.R))
	}
	if len(mp.PendingFlowers) != 1 {
		t.Errorf("pending flowers should be 1 %d", len(mp.PendingFlowers))
	}
	if len(mp.HE) != 0 {
		t.Errorf("east seg size should be 0 %d", len(mp.HE))
	}
	if len(mp.HS) != 0 {
		t.Errorf("south seg size should be 0 %d", len(mp.HS))
	}
	if len(mp.HW) != 0 {
		t.Errorf("west seg size should be 0 %d", len(mp.HW))
	}
	if len(mp.HN) != 0 {
		t.Errorf("north seg size should be 0 %d", len(mp.HN))
	}
	if len(mp.G) != 0 {
		t.Errorf("green seg size should be 0 %d", len(mp.G))
	}
	if len(mp.W) != 0 {
		t.Errorf("white seg size should be 0 %d", len(mp.W))
	}
	if mp.TC[TC_B] != 3 {
		t.Errorf("bamboo seg size should be 3 %d", mp.TC[TC_B])
	}
	if mp.TC[TC_C] != 3 {
		t.Errorf("characters seg size should be 3 %d", mp.TC[TC_C])
	}
	if mp.TC[TC_D] != 4 {
		t.Errorf("dots seg size should be 4 %d", mp.TC[TC_D])
	}
	if mp.TC[TC_H] != 2 {
		t.Errorf("hornor seg size should be 4 %d", mp.TC[TC_H])
	}
	if mp.LD != mj.FromS(mj.RED).Seq {
		t.Errorf("last draw should be %d , %d", mj.FromS(mj.RED).Seq, mp.LD)
	}
	if mp.TN {
		t.Errorf("hand state should no setting as false %v", mp.TN)
	}
	if mp.SystemId != 0 {
		t.Errorf("seat should no player as 0 %v", mp.SystemId)
	}
}
func TestPlayerChow(t *testing.T) {
	core.CreateTestLog()
	mp := NewPlayer(0, true, &SampleCallback{})
	mp.AppendForTest(mj.FromS(mj.BAMBOO1), false)
	mp.AppendForTest(mj.FromS(mj.BAMBOO2), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER1), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER3), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER7), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER8), false)
	err := mp.Chow(mj.FromS(mj.BAMBOO3), mj.FromS(mj.BAMBOO1), mj.FromS(mj.BAMBOO2))
	if err != nil {
		t.Errorf("should be a chow %s", err.Error())
	}
	formed := false
	for i := range mp.Formed {
		meld := mp.Formed[i]
		if meld.Name() == "B1.B2.B3" {
			formed = true
			break
		}
	}
	if !formed {
		t.Errorf("should be a chow meld %s", "B1.B2.B3")
	}
	bz := len(mp.B)
	if bz > 0 {
		t.Errorf("should be no bamboo %d", bz)
	}
	hz := mp.TileSize()
	if hz != 4 {
		t.Errorf("hand size should be 4%d", hz)
	}
	err = mp.Chow(mj.FromS(mj.CHARACTER2), mj.FromS(mj.CHARACTER1), mj.FromS(mj.CHARACTER3))
	if err != nil {
		t.Errorf("should be a chow %s", err.Error())
	}
	formed = false
	for i := range mp.Formed {
		meld := mp.Formed[i]
		if meld.Name() == "C1.C2.C3" {
			formed = true
			break
		}
	}
	if !formed {
		t.Errorf("should be a chow meld %s", "C1.C2.C3")
	}
	cz := len(mp.C)
	if cz != 2 {
		t.Errorf("should be 2 characters %d", cz)
	}
	hz = mp.TileSize()
	if hz != 2 {
		t.Errorf("hand size should be 2%d", hz)
	}

	err = mp.Chow(mj.FromS(mj.CHARACTER9), mj.FromS(mj.CHARACTER7), mj.FromS(mj.CHARACTER8))
	if err != nil {
		t.Errorf("should be a chow %s", err.Error())
	}
	formed = false
	for i := range mp.Formed {
		meld := mp.Formed[i]
		if meld.Name() == "C7.C8.C9" {
			formed = true
			break
		}
	}
	if !formed {
		t.Errorf("should be a chow meld %s", "C7.C8.C9")
	}
	cz = len(mp.C)
	if cz != 0 {
		t.Errorf("should be 0 characters %d", cz)
	}
	hz = mp.TileSize()
	if hz != 0 {
		t.Errorf("hand size should be 0%d", hz)
	}
	fz := len(mp.Formed)
	if fz != 3 {
		t.Errorf("formed size should be 3%d", fz)
	}
}

func TestPlayerPung(t *testing.T) {
	core.CreateTestLog()
	mp := NewPlayer(0, true, &SampleCallback{})
	mp.AppendForTest(mj.FromS(mj.EAST), false)
	mp.AppendForTest(mj.FromS(mj.EAST), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER1), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER1), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER3), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER3), false)
	err := mp.Pung(mj.FromS(mj.EAST))
	if err != nil {
		t.Errorf("should be a pung %s", err.Error())
	}
	formed := false
	for i := range mp.Formed {
		meld := mp.Formed[i]
		if meld.Name() == "H31.H31.H31" {
			formed = true
			break
		}
	}
	if !formed {
		t.Errorf("should be a pung meld %s", "H31.H31.H31")
	}
	err = mp.Pung(mj.FromS(mj.CHARACTER1))
	if err != nil {
		t.Errorf("should be a pung %s", err.Error())
	}
	formed = false
	for i := range mp.Formed {
		meld := mp.Formed[i]
		if meld.Name() == "C1.C1.C1" {
			formed = true
			break
		}
	}
	if !formed {
		t.Errorf("should be a pung meld %s", "C1.C1.C1")
	}
	err = mp.Pung(mj.FromS(mj.CHARACTER3))
	if err != nil {
		t.Errorf("should be a pung %s", err.Error())
	}
	formed = false
	for i := range mp.Formed {
		meld := mp.Formed[i]
		if meld.Name() == "C3.C3.C3" {
			formed = true
			break
		}
	}
	if !formed {
		t.Errorf("should be a pung meld %s", "C3.C3.C3")
	}
	hz := mp.TileSize()
	if hz != 0 {
		t.Errorf("hand size should be 0%d", hz)
	}
	fz := len(mp.Formed)
	if fz != 3 {
		t.Errorf("formed size should be 3%d", fz)
	}
}

func TestPlayerKong(t *testing.T) {
	core.CreateTestLog()
	mp := NewPlayer(0, true, &SampleCallback{})
	mp.AppendForTest(mj.FromS(mj.EAST), false)
	mp.AppendForTest(mj.FromS(mj.EAST), false)
	mp.AppendForTest(mj.FromS(mj.EAST), false)
	mp.AppendForTest(mj.FromS(mj.EAST), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER3), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER3), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER3), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER2), false)
	//mp.AppendForTest(mj.FromS(mj.CHARACTER7), false)
	//mp.AppendForTest(mj.FromS(mj.CHARACTER7), false)
	kt, err := mp.validateKong(mj.FromS(mj.EAST).Seq)
	if err != nil {
		t.Errorf("should be a kong %s", err.Error())
	}
	if kt.Tye != K_CONCEALED {
		t.Errorf("should be a concealed kong %d", kt.Tye)
	}
	err = mp.Kong(mj.FromS(mj.EAST)) //concealed kong
	if err != nil {
		t.Errorf("should be a kong %s", err.Error())
	}
	formed := false
	for i := range mp.Formed {
		meld := mp.Formed[i]
		if meld.Name() == "H31.H31.H31.H31" && meld.Concealed {
			formed = true
			break
		}
	}
	if !formed {
		t.Errorf("should be a kong meld %s", "H31.H31.H31,H31")
	}
	err = mp.Kong(mj.FromS(mj.CHARACTER3)) //pung kong as exposed kong
	if err != nil {
		t.Errorf("should be a kong %s", err.Error())
	}
	formed = false
	for i := range mp.Formed {
		meld := mp.Formed[i]
		if meld.Name() == "C3.C3.C3.C3" && (!meld.Concealed) {
			formed = true
			break
		}
	}
	if !formed {
		t.Errorf("should be a pung meld %s", "C3.C3.C3,C3")
	}
	p2 := mj.Meld{}
	p2.Tiles = []mj.Tile{mj.FromS(mj.CHARACTER1), mj.FromS(mj.CHARACTER1), mj.FromS(mj.CHARACTER1)}
	mp.Formed = append(mp.Formed, p2)
	err = mp.Kong(mj.FromS(mj.CHARACTER1)) //draw to kong from meld pung
	if err != nil {
		t.Errorf("should be a kong %s", err.Error())
	}
	formed = false
	for i := range mp.Formed {
		meld := mp.Formed[i]
		if meld.Name() == "C1.C1.C1.C1" && (!meld.Concealed) {
			formed = true
			break
		}
	}
	if !formed {
		t.Errorf("should be a pung meld %s", "C1.C1.C1,C1")
	}
	hz := mp.TileSize()
	if hz != 1 {
		t.Errorf("hand size should be 1%d", hz)
	}
	fz := len(mp.Formed)
	if fz != 3 {
		t.Errorf("formed size should be 3%d", fz)
	}
	err = mp.Kong(mj.FromS(mj.CHARACTER2)) //no kong
	if err == nil {
		t.Errorf("should be no kong")
	}
}
