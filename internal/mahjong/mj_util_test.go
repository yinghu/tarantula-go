package mahjong

import (
	"slices"
	"testing"
)




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

func TestSuit(t *testing.T) {
	d1 := Tile{}
	d1.From(DOTS1)
	if d1.Rank != 1 {
		t.Errorf("value should be 1 %d", d1.Rank)
	}
	if d1.Suit != "D" {
		t.Errorf("SUIT should be D %s", d1.Suit)
	}
	c9 := Tile{}
	c9.From(CHARACTER9)
	if c9.Rank != 9 {
		t.Errorf("value should be 9 %d", c9.Rank)
	}
	if c9.Suit != "C" {
		t.Errorf("SUIT should be C %s", c9.Suit)
	}

	f3 := Tile{}
	f3.From(F_CHRYSANTHEMUM)
	if f3.Rank != 103 {
		t.Errorf("value should be 103 %d", f3.Rank)
	}
	if f3.Suit != "F" {
		t.Errorf("SUIT should be F %s", f3.Suit)
	}

}

func TestShuffle(t *testing.T) {
	s := Deck{}
	s.New()
	if len(s.Stack) != s.Size {
		t.Errorf("deck size error %d", len(s.Stack))
	}
	s.Shuffle()
	if s.header != 0 {
		t.Errorf("deck header error %d", s.header)
	}
	if s.tail == s.Size {
		t.Errorf("deck tail error %d", s.tail)
	}
	_, err := s.Draw()
	if err != nil {
		t.Errorf("first draw erro %s", err.Error())
	}
	if s.header != 1 {
		t.Errorf("first draw erro %d", s.header)
	}
	tail := s.tail
	_, err = s.Kong()
	if err != nil {
		t.Errorf("first kong erro %s", err.Error())
	}
	if s.tail != tail-1 {
		t.Errorf("first kong erro %d", s.tail)
	}
	s.tail = 1
	_, err = s.Draw()
	if err != nil {
		t.Errorf("should be last draw %s", err.Error())
	}
	noTile, err := s.Draw()
	if err == nil {
		t.Errorf("should be no more draw %s", noTile.Suit)
	}
	s.header = 100
	s.tail = 100
	_, err = s.Kong()
	if err != nil {
		t.Errorf("should be last knog %s", err.Error())
	}
	noKnog, err := s.Kong()
	if err == nil {
		t.Errorf("should be no more know %s", noKnog.Suit)
	}
}






func TestHandStart(t *testing.T) {
	s := Deck{}
	s.New()
	s.Shuffle()
	h := Hand{}
	h.New()
	err := h.Draw(&s)
	if err != nil {
		t.Errorf("should be error")
	}
	err = h.Knog(&s)
	if err != nil {
		t.Errorf("should be error")
	}

}

