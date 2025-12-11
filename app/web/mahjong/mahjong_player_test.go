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
	mds := mp.CheckDiscard(1, mj.FromS(mj.CHARACTER1), false)
	for i := range mds {
		fmt.Printf("Kong %s\n", mds[i].Name())
	}
	fmt.Printf("pending kong %v\n", mp.PendingKongs)
	err = mp.Pung(mj.FromS(mj.BAMBOO1))
	if err != nil {
		t.Errorf("should be a ping %s", err.Error())
	}
	err = mp.AppendForTest(mj.FromS(mj.BAMBOO1),false)
	if err != nil {
		t.Errorf("should be no error %s", err.Error())
	}
	fmt.Printf("pending kong %v\n", mp.PendingKongs)
	
}
