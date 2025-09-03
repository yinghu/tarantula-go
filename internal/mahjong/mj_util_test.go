package mahjong

import (
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
	d5 := Tile{Suit: "C", Rank: 5} //FIRST
	d6 := Tile{Suit: "D", Rank: 6}
	d7 := Tile{Suit: "D", Rank: 7} //LAST
	d8 := Tile{Suit: "D", Rank: 1}

	ts := []Tile{d7, d1, d6, d2, d5, d8}
	h := Hand{}
	h.New()
	h.Tiles = append(h.Tiles, ts...)
	slices.SortFunc(h.Tiles, cmp)
	//fmt.Printf("%v\n", h.Tiles)
	if h.Tiles[0] != d5 {
		t.Errorf("first should be d 5")
	}
	if h.Tiles[5] != d7 {
		t.Errorf("last should be d 7")
	}
}

func TestSlice(t *testing.T) {
	ts := make([]Tile, 0)
	ts = append(ts, Tile{Suit: "B", Rank: 1})
	ts = append(ts, Tile{Suit: "B", Rank: 2})
	ts = append(ts, Tile{Suit: "B", Rank: 3})
	sz := len(ts)
	if sz != 3 {
		t.Errorf("size should be 3 %d", sz)
	}
	ts = ts[1:]
	sz = len(ts)
	if sz != 2 {
		t.Errorf("size should be 2 %d", sz)
	}

	d1 := Tile{Suit: "D", Rank: 1}
	d2 := Tile{Suit: "D", Rank: 2}
	d3 := Tile{Suit: "D", Rank: 3}
	xs := make([]Tile, 0)
	xs = slices.Insert(xs, 0, d1)
	xs = slices.Insert(xs, 0, d2)
	xs = slices.Insert(xs, 0, d3)
	xz := len(xs)
	if xz != 3 {
		t.Errorf("size should be 3 %d", xz)
	}
	if xs[0] != d3 {
		t.Errorf("first should be d 3")
	}
	if xs[1] != d2 {
		t.Errorf("second should be d 2")
	}
	if xs[2] != d1 {
		t.Errorf("third should be d 1")
	}
}

func TestStack(t *testing.T) {
	h := Hand{}
	h.New()
	ts1 := h.NewTileSet(FOUR_SET)
	ts2 := h.NewTileSet(THREE_SET)
	ts3 := h.NewTileSet(SEQ_SET)
	ts4 := h.NewTileSet(TWO_SET)
	h.Push(ts1)
	h.Push(ts2)
	h.Push(ts3)
	h.Push(ts4)
	if h.StackSize() != 4 {
		t.Errorf("stack side should be 4 %d", h.StackSize())
	}

}

func TestNextTiles(t *testing.T) {
	h := Hand{}
	h.New()
	b1 := Tile{Suit: "B", Rank: 1}
	b2 := Tile{Suit: "B", Rank: 2}
	b3 := Tile{Suit: "B", Rank: 3}

	b4 := Tile{Suit: "B", Rank: 2}
	b5 := Tile{Suit: "B", Rank: 3}
	b6 := Tile{Suit: "B", Rank: 4}

	b7 := Tile{Suit: "B", Rank: 3}
	b8 := Tile{Suit: "B", Rank: 4}
	b9 := Tile{Suit: "B", Rank: 5}

	c7 := Tile{Suit: "B", Rank: 4}
	c8 := Tile{Suit: "B", Rank: 5}
	c9 := Tile{Suit: "B", Rank: 6}

	p1 := Tile{Suit: "B", Rank: 1}
	p2 := Tile{Suit: "B", Rank: 1}
	h.Tiles = append(h.Tiles, b1, b2, b3, b4, b5, b6, b7, b8, b9, c7, c8, c9, p1, p2)
	next := h.NextTile()
	if next != b1 {
		t.Errorf("next tile should be b1 %v", b1)
	}
	h.PopTile()
	h.PopTile()
	h.PopTile()
	next = h.NextTile()
	if next != b4 {
		t.Errorf("next tile should be b4 %v", b4)
	}
	nexts := h.NextTiles(100)
	if len(nexts) != 11 {
		t.Errorf("remaining list size should be 11 %d",len(h.Tiles))
	}

	nexts = h.NextTiles(3)
	if len(nexts) != 3 {
		t.Errorf("next size should be 3 %d",len(nexts))
	}
	if nexts[0] != b4{
		t.Errorf("first should be b4 %v",nexts[0])
	}
	if nexts[1] != b5{
		t.Errorf("first should be b5 %v",nexts[0])
	}
	if nexts[2] != b6{
		t.Errorf("first should be b6 %v",nexts[0])
	}

}
