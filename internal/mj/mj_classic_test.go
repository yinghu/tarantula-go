package mj

import (
	"testing"
)

func TestClassic(t *testing.T) {

	b1 := FromS(BAMBOO5)
	b2 := FromS(BAMBOO6)
	b3 := FromS(BAMBOO7)

	b4 := FromS(BAMBOO7)
	b5 := FromS(BAMBOO8)
	b6 := FromS(BAMBOO9)

	b7 := FromS(CHARACTER3)
	b8 := FromS(CHARACTER4)
	b9 := FromS(CHARACTER5)

	c7 := FromS(CHARACTER6)
	c8 := FromS(CHARACTER6)

	c9 := FromS(RED)

	p1 := FromS(RED)
	p2 := FromS(RED)

	tiles := []Tile{b1, b2, b3, b4, b5, b6, b7, b8, b9, c7, c8, c9, p1, p2}
	h := Hand{}
	h.New()
	h.Tiles = append(h.Tiles, tiles...)
	cm := ClassicMahjong{}
	cm.New()
	claimed := cm.Mahjong(&h)
	if !claimed {
		t.Errorf("should be claimed %v", claimed)
	}
	//B1,B2,B3,B2,B3,B4,B3,B4,B5,B4,B5,B6,B1,B1
}

func TestClassic1(t *testing.T) {

	b1 := FromS(BAMBOO1)
	b2 := FromS(BAMBOO1)
	b3 := FromS(BAMBOO2)

	b4 := FromS(BAMBOO3)
	b5 := FromS(BAMBOO4)
	b6 := FromS(BAMBOO6)

	b7 := FromS(BAMBOO6)
	b8 := FromS(BAMBOO7)
	b9 := FromS(BAMBOO7)

	c7 := FromS(BAMBOO8)
	c8 := FromS(BAMBOO8)

	c9 := FromS(DOTS5)

	p1 := FromS(DOTS6)
	p2 := FromS(DOTS7)

	tiles := []Tile{b1, b2, b3, b4, b5, b6, b7, b8, b9, c7, c8, c9, p1, p2}
	h := Hand{}
	h.New()
	h.Tiles = append(h.Tiles, tiles...)
	//slices.SortFunc(h.Tiles, cmp)
	cm := ClassicMahjong{}
	cm.New()
	claimed := cm.Mahjong(&h)
	if !claimed {
		t.Errorf("should be claimed %v", claimed)
	}
	//B1,B2,B3,B2,B3,B4,B3,B4,B5,B4,B5,B6,B1,B1
}
