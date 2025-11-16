package mj

import (
	"testing"
)

func TestClassic(t *testing.T) {

	b1 := NewTile(BAMBOO5)
	b2 := NewTile(BAMBOO6)
	b3 := NewTile(BAMBOO7)

	b4 := NewTile(BAMBOO7)
	b5 := NewTile(BAMBOO8)
	b6 := NewTile(BAMBOO9)

	b7 := NewTile(CHARACTER3)
	b8 := NewTile(CHARACTER4)
	b9 := NewTile(CHARACTER5)

	c7 := NewTile(CHARACTER6)
	c8 := NewTile(CHARACTER6)

	c9 := NewTile(RED)

	p1 := NewTile(RED)
	p2 := NewTile(RED)

	tiles := []Tile{b1, b2, b3, b4, b5, b6, b7, b8, b9, c7, c8, c9, p1, p2}
	h := Hand{}
	h.New()
	h.Tiles = append(h.Tiles, tiles...)
	cm := ClassicMahjong{}
	cm.New()
	claimed := cm.Mahjong(&h)
	if !claimed{
		t.Errorf("should be claimed %v",claimed)
	}
	//B1,B2,B3,B2,B3,B4,B3,B4,B5,B4,B5,B6,B1,B1
}
