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

