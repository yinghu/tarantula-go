package mahjong

import (
	"testing"
)

func TestHandEval1(t *testing.T) {
	h := Hand{}
	h.New()
	d1 := Tile{Suit: "D", Rank: 1}
	d2 := Tile{Suit: "D", Rank: 2}
	d3 := Tile{Suit: "D", Rank: 3}

	c4 := Tile{Suit: "C", Rank: 4}
	c5 := Tile{Suit: "C", Rank: 5}
	c6 := Tile{Suit: "C", Rank: 6}

	b1 := Tile{Suit: "B", Rank: 1}
	b2 := Tile{Suit: "B", Rank: 2}
	b3 := Tile{Suit: "B", Rank: 3}

	b7 := Tile{Suit: "B", Rank: 7}
	b8 := Tile{Suit: "B", Rank: 8}
	b9 := Tile{Suit: "B", Rank: 9}

	p1 := Tile{Suit: "D", Rank: 2}
	p2 := Tile{Suit: "D", Rank: 2}
	h.Tiles = append(h.Tiles, d1, d2, d3, c4, c5, c6, b1, b2, b3, b7, b8, b9, p1, p2)
	if h.TileSize() != 14 {
		t.Errorf("hand size should be 14 %d", h.TileSize())
	}
	matched := h.Mahjong()
	if !matched {
		t.Errorf("hand should be a match %v", matched)
	}
}

func TestHandEval2(t *testing.T) {
	h := Hand{}
	h.New()
	d1 := Tile{Suit: "D", Rank: 1}
	d2 := Tile{Suit: "D", Rank: 2}
	d3 := Tile{Suit: "D", Rank: 3}

	c4 := Tile{Suit: "C", Rank: 4}
	c5 := Tile{Suit: "C", Rank: 5}
	c6 := Tile{Suit: "C", Rank: 6}

	b1 := Tile{Suit: "B", Rank: 1}
	b2 := Tile{Suit: "B", Rank: 2}
	b3 := Tile{Suit: "B", Rank: 3}

	b7 := Tile{Suit: "B", Rank: 7}
	b8 := Tile{Suit: "B", Rank: 8}
	b9 := Tile{Suit: "B", Rank: 9}

	p1 := Tile{Suit: "H", Rank: 2}
	p2 := Tile{Suit: "H", Rank: 2}
	h.Tiles = append(h.Tiles, d1, d2, d3, c4, c5, c6, b1, b2, b3, b7, b8, b9, p1, p2)
	if h.TileSize() != 14 {
		t.Errorf("hand size should be 14 %d", h.TileSize())
	}
	matched := h.Mahjong()
	if !matched {
		t.Errorf("hand should be a match %v", matched)
	}
}

func TestHandEval3(t *testing.T) {
	h := Hand{}
	h.New()
	d1 := Tile{Suit: "D", Rank: 1}
	d2 := Tile{Suit: "D", Rank: 2}
	d3 := Tile{Suit: "D", Rank: 3}

	c4 := Tile{Suit: "D", Rank: 2}
	c5 := Tile{Suit: "D", Rank: 3}
	c6 := Tile{Suit: "D", Rank: 4}

	b1 := Tile{Suit: "B", Rank: 1}
	b2 := Tile{Suit: "B", Rank: 2}
	b3 := Tile{Suit: "B", Rank: 3}

	b7 := Tile{Suit: "B", Rank: 7}
	b8 := Tile{Suit: "B", Rank: 8}
	b9 := Tile{Suit: "B", Rank: 9}

	p1 := Tile{Suit: "H", Rank: 2}
	p2 := Tile{Suit: "H", Rank: 2}
	h.Tiles = append(h.Tiles, d1, d2, d3, c4, c5, c6, b1, b2, b3, b7, b8, b9, p1, p2)
	if h.TileSize() != 14 {
		t.Errorf("hand size should be 14 %d", h.TileSize())
	}
	matched := h.Mahjong()
	if !matched {
		t.Errorf("hand should be a match %v", matched)
	}
}