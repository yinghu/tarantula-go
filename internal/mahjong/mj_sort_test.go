package mahjong

import (
	"fmt"
	"slices"
	"testing"
)

func TestTileRankSort(t *testing.T) {
	d1 := Tile{Suit: "D", Rank: 9}
	d2 := Tile{Suit: "D", Rank: 2}
	d5 := Tile{Suit: "D", Rank: 5}
	d6 := Tile{Suit: "D", Rank: 6}
	d7 := Tile{Suit: "D", Rank: 7}
	d8 := Tile{Suit: "D", Rank: 1}

	ts := []Tile{d7, d1, d6, d2, d5, d8}
	h := Hand{}
	h.New()
	h.Tiles = append(h.Tiles, ts...)
	slices.SortFunc(h.Tiles, cmp)
	if h.Tiles[0] != d8 {
		t.Errorf("first should be d 8")
	}
	if h.Tiles[5] != d1 {
		t.Errorf("last should be d 1")
	}
}

func TestTileMixSort(t *testing.T) {
	d1 := Tile{Suit: "C", Rank: 9}
	d2 := Tile{Suit: "D", Rank: 2}
	d5 := Tile{Suit: "C", Rank: 5}//FIRST
	d6 := Tile{Suit: "D", Rank: 6}
	d7 := Tile{Suit: "D", Rank: 7} //LAST
	d8 := Tile{Suit: "D", Rank: 1}

	ts := []Tile{d7, d1, d6, d2, d5, d8}
	h := Hand{}
	h.New()
	h.Tiles = append(h.Tiles, ts...)
	slices.SortFunc(h.Tiles, cmp)
	fmt.Printf("%v\n",h.Tiles)
	if h.Tiles[0] != d5 {
		t.Errorf("first should be d 5")
	}
	if h.Tiles[5] != d7 {
		t.Errorf("last should be d 7")
	}
}
