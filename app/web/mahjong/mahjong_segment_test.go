package main

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/mj"
)

func TestMahjongSegment(t *testing.T) {
	core.CreateTestLog()
	tiles := make([]mj.Tile, 0)
	tiles = append(tiles, mj.FromS(mj.BAMBOO1))
	tiles = append(tiles, mj.FromS(mj.BAMBOO2))
	tiles = append(tiles, mj.FromS(mj.BAMBOO3))
	tiles = append(tiles, mj.FromS(mj.BAMBOO4))
	tiles = append(tiles, mj.FromS(mj.BAMBOO5))
	tiles = append(tiles, mj.FromS(mj.BAMBOO3))
	tiles = append(tiles, mj.FromS(mj.BAMBOO3))

	tiles = append(tiles, mj.FromS(mj.BAMBOO6))
	tiles = append(tiles, mj.FromS(mj.BAMBOO6))
	tiles = append(tiles, mj.FromS(mj.BAMBOO6))
	seg := HandSegmenet{}
	seg.From(tiles)
	pn := seg.AfterPair()
	if pn != 2 {
		t.Errorf("pair number should 2 %d", pn)
	}
	b6 := seg.Index[mj.FromS(mj.BAMBOO6).Seq]
	if b6.Count-b6.Used != 1 {
		t.Errorf("b6 used number should 2 %d", b6.Used)
	}
	b3 := seg.Index[mj.FromS(mj.BAMBOO3).Seq]
	if b3.Count-b3.Used != 1 {
		t.Errorf("b3 used number should 2 %d", b3.Used)
	}
	pc := seg.AfterChow()
	if pc != 2 {
		t.Errorf("chow number should 2 %d", pc)
	}
	for i := range seg.Index {
		c := seg.Index[i]
		fmt.Printf("Usage : %d >> %d >> %d\n", i, c.Count, c.Used)
	}

}
