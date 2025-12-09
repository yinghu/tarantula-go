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
	seq := seg.discard(tiles)
	fmt.Printf("Discard %v\n",mj.FromQ(seq))
	seg.From(tiles)
	seg.AfterKong()
	seg.AfterPung()
	seg.AfterChow()
	for s := range seg.Index {
		v := seg.Index[s]
		fmt.Printf("Used %d %d %d\n", s, v.Count,v.Used)
	}
	seg.Reset()
	seg.AfterChow()
	seg.AfterKong()
	seg.AfterPung()
	for s := range seg.Index {
		v := seg.Index[s]
		fmt.Printf("Used %d %d %d\n", s, v.Count,v.Used)
	}


}
