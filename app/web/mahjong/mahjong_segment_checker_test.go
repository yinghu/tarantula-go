package main

import (
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
	seg := HandSegmenet{Used: make(map[int]int)}
	seg.Listener = &seg
	seg.From(tiles)
	seg.Chow()
	seg.Pung()
	seg.Kong()

}
