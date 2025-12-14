package main

import (
	"fmt"
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
	mp := NewPlayer(0, true, &SampleCallback{})
	if !mp.Auto {
		t.Errorf("default should be auto %v", mp.Auto)
	}
	if mp.FlowerExcluded {
		t.Errorf("flower should be included %v", mp.FlowerExcluded)
	}
	mp.AppendForTest(mj.FromS(mj.BAMBOO1), false)
	mp.AppendForTest(mj.FromS(mj.BAMBOO1), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER1), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER1), false)
	mp.AppendForTest(mj.FromS(mj.CHARACTER1), false)
	mp.AppendForTest(mj.FromS(mj.DOTS1), false)
	mp.AppendForTest(mj.FromS(mj.DOTS1), false)
	mp.AppendForTest(mj.FromS(mj.DOTS1), false)
	mp.AppendForTest(mj.FromS(mj.DOTS1), false)
	mp.AppendForTest(mj.FromS(mj.RED), false)
	mp.AppendForTest(mj.FromS(mj.F_AUTUMN), false)
	fmt.Printf("hand %v\n", mp.Hand.Tiles)
	fmt.Printf("bamboo %v\n", mp.B)
	fmt.Printf("character %v\n", mp.C)
	fmt.Printf("dots %v\n", mp.D)
	fmt.Printf("red %v\n", mp.R)
	fmt.Printf("distribution %v\n", mp.TC)
	fmt.Printf("pending kong %v\n", mp.PendingKongs)
	fmt.Printf("last draw %d\n", mp.LD)
	err := mp.validateKong(mj.FromS(mj.DOTS1).Seq)
	fmt.Printf("pending kong %v\n", mp.PendingKongs)
	if err != nil {
		t.Errorf("should be a kong %s", err.Error())
	}
	err = mp.Kong(mj.FromS(mj.DOTS1))
	if err != nil {
		t.Errorf("should be a kong %s", err.Error())
	}
	fmt.Printf("hand %v\n", mp.Hand.Tiles)
	mds := mp.CheckDiscard(1, mj.FromS(mj.CHARACTER1), false)
	for i := range mds {
		fmt.Printf("Kong %s\n", mds[i].Name())
	}
	fmt.Printf("pending kong %v\n", mp.PendingKongs)
	err = mp.Pung(mj.FromS(mj.BAMBOO1))
	if err != nil {
		t.Errorf("should be a ping %s", err.Error())
	}
	err = mp.AppendForTest(mj.FromS(mj.BAMBOO1), false)
	if err != nil {
		t.Errorf("should be no error %s", err.Error())
	}
	fmt.Printf("pending kongs %v\n", mp.PendingKongs)
	fmt.Printf("pending flowers %v\n", mp.PendingFlowers)
	seg := make([]mj.Tile, 0)
	seg = append(seg, mj.FromS(mj.BAMBOO7))
	seg = append(seg, mj.FromS(mj.BAMBOO9))
	seg = append(seg, mj.FromS(mj.BAMBOO9))
	seg = append(seg, mj.FromS(mj.BAMBOO7))
	seg = append(seg, mj.FromS(mj.BAMBOO8))
	mp.checkKong(seg, false)
	fmt.Printf("pending kongs %v\n", mp.PendingKongs)
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
}

func TestPlayerKong(t *testing.T) {

}
