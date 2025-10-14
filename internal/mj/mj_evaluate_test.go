package mj

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

func TestHandEval4(t *testing.T) {
	h := Hand{}
	h.New()
	d1 := Tile{Suit: "B", Rank: 1}
	d2 := Tile{Suit: "B", Rank: 2}
	d3 := Tile{Suit: "B", Rank: 3}

	c4 := Tile{Suit: "B", Rank: 2}
	c5 := Tile{Suit: "B", Rank: 3}
	c6 := Tile{Suit: "B", Rank: 4}

	b1 := Tile{Suit: "D", Rank: 1}
	b2 := Tile{Suit: "D", Rank: 2}
	b3 := Tile{Suit: "D", Rank: 3}

	b7 := Tile{Suit: "C", Rank: 7}
	b8 := Tile{Suit: "C", Rank: 8}
	b9 := Tile{Suit: "C", Rank: 9}

	p1 := Tile{Suit: "C", Rank: 8}
	p2 := Tile{Suit: "C", Rank: 8}
	h.Tiles = append(h.Tiles, d1, d2, d3, c4, c5, c6, b1, b2, b3, b7, b8, b9, p1, p2)
	if h.TileSize() != 14 {
		t.Errorf("hand size should be 14 %d", h.TileSize())
	}
	matched := h.Mahjong()
	if !matched {
		t.Errorf("hand should be a match %v", matched)
	}
}

func TestHandEval5(t *testing.T) {
	h := Hand{}
	h.New()
	b1 := Tile{Suit: "B", Rank: 1}
	b2 := Tile{Suit: "B", Rank: 2}
	b3 := Tile{Suit: "B", Rank: 3}

	b4 := Tile{Suit: "B", Rank: 2}
	b5 := Tile{Suit: "B", Rank: 3}
	b6 := Tile{Suit: "B", Rank: 4}

	b7 := Tile{Suit: "B", Rank: 2}
	b8 := Tile{Suit: "B", Rank: 3}
	b9 := Tile{Suit: "B", Rank: 4}

	c7 := Tile{Suit: "C", Rank: 7}
	c8 := Tile{Suit: "C", Rank: 8}
	c9 := Tile{Suit: "C", Rank: 9}

	p1 := Tile{Suit: "C", Rank: 8}
	p2 := Tile{Suit: "C", Rank: 8}
	h.Tiles = append(h.Tiles, b1, b2, b3, b4, b5, b6, b7, b8, b9, c7, c8, c9, p1, p2)
	if h.TileSize() != 14 {
		t.Errorf("hand size should be 14 %d", h.TileSize())
	}
	matched := h.Mahjong()
	if !matched {
		t.Errorf("hand should be a match %v", matched)
	}
}

func TestHandEval6(t *testing.T) {
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

	p1 := Tile{Suit: "B", Rank: 6}
	p2 := Tile{Suit: "B", Rank: 6}
	h.Tiles = append(h.Tiles, b1, b2, b3, b4, b5, b6, b7, b8, b9, c7, c8, c9, p1, p2)
	if h.TileSize() != 14 {
		t.Errorf("hand size should be 14 %d", h.TileSize())
	}
	matched := h.Mahjong()
	if !matched {
		t.Errorf("hand should be a match %v", matched)
	}
}

func TestHandEval8(t *testing.T) {
	h := Hand{}
	h.New()
	b1 := Tile{Suit: "B", Rank: 1}
	b2 := Tile{Suit: "B", Rank: 1}
	b3 := Tile{Suit: "B", Rank: 1}

	b4 := Tile{Suit: "B", Rank: 3}
	b5 := Tile{Suit: "B", Rank: 3}
	b6 := Tile{Suit: "B", Rank: 3}

	b7 := Tile{Suit: "B", Rank: 3}
	b8 := Tile{Suit: "B", Rank: 4}
	b9 := Tile{Suit: "B", Rank: 5}

	c7 := Tile{Suit: "B", Rank: 4}
	c8 := Tile{Suit: "B", Rank: 5}
	c9 := Tile{Suit: "B", Rank: 6}

	p1 := Tile{Suit: "B", Rank: 7}
	p2 := Tile{Suit: "B", Rank: 7}
	h.Tiles = append(h.Tiles, b1, b2, b3, b4, b5, b6, b7, b8, b9, c7, c8, c9, p1, p2)
	if h.TileSize() != 14 {
		t.Errorf("hand size should be 14 %d", h.TileSize())
	}
	matched := h.Mahjong()
	if !matched {
		t.Errorf("hand should be a match %v", matched)
	}
}

func TestHandEval9(t *testing.T) {
	h := Hand{}
	h.New()
	b1 := Tile{Suit: "B", Rank: 1}
	b2 := Tile{Suit: "B", Rank: 3}
	b3 := Tile{Suit: "B", Rank: 5}

	b4 := Tile{Suit: "C", Rank: 3}
	b5 := Tile{Suit: "C", Rank: 6}
	b6 := Tile{Suit: "C", Rank: 9}

	b7 := Tile{Suit: "D", Rank: 2}
	b8 := Tile{Suit: "D", Rank: 5}
	b9 := Tile{Suit: "D", Rank: 8}

	c7 := Tile{Suit: "H", Rank: 1}
	c8 := Tile{Suit: "H", Rank: 4}
	c9 := Tile{Suit: "H", Rank: 5}

	p1 := Tile{Suit: "H", Rank: 2}
	p2 := Tile{Suit: "H", Rank: 9}
	h.Tiles = append(h.Tiles, b1, b2, b3, b4, b5, b6, b7, b8, b9, c7, c8, c9, p1, p2)
	if h.TileSize() != 14 {
		t.Errorf("hand size should be 14 %d", h.TileSize())
	}
	matched := h.Mahjong()
	if !matched {
		t.Errorf("hand should be a match %v", matched)
	}
}

func TestHandEval10(t *testing.T) {
	h := Hand{}
	h.New()
	b1 := Tile{Suit: "B", Rank: 1}
	b2 := Tile{Suit: "B", Rank: 2}
	b3 := Tile{Suit: "B", Rank: 3}

	b4 := Tile{Suit: "C", Rank: 3}
	b5 := Tile{Suit: "C", Rank: 3}
	b6 := Tile{Suit: "C", Rank: 3}

	b7 := Tile{Suit: "D", Rank: 2}
	b8 := Tile{Suit: "D", Rank: 2}
	b9 := Tile{Suit: "D", Rank: 2}

	c7 := Tile{Suit: "B", Rank: 3}
	c8 := Tile{Suit: "B", Rank: 4}
	c9 := Tile{Suit: "B", Rank: 5}

	p1 := Tile{Suit: "H", Rank: 7}
	p2 := Tile{Suit: "H", Rank: 7}
	h.Tiles = append(h.Tiles, b1, b2, b3, b4, b5, b6, b7, b8, b9, c7, c8, c9, p1, p2)
	if h.TileSize() != 14 {
		t.Errorf("hand size should be 14 %d", h.TileSize())
	}
	matched := h.Mahjong()
	if !matched {
		t.Errorf("hand should be a match %v", matched)
	}
}

func TestHandEval11(t *testing.T) {

	b1 := NewTile(BAMBOO1)
	b2 := NewTile(BAMBOO2)
	b3 := NewTile(BAMBOO3)

	b4 := NewTile(BAMBOO2)
	b5 := NewTile(BAMBOO3)
	b6 := NewTile(BAMBOO4)

	b7 := NewTile(BAMBOO3)
	b8 := NewTile(BAMBOO4)
	b9 := NewTile(BAMBOO5)

	c7 := NewTile(BAMBOO4)
	c8 := NewTile(BAMBOO5)
	c9 := NewTile(BAMBOO6)

	p1 := NewTile(BAMBOO1)
	p2 := NewTile(BAMBOO1)

	tiles := []Tile{b1, b2, b3, b4, b5, b6, b7, b8, b9, c7, c8, c9, p1, p2}
	h := Hand{}
	h.New()
	h.Tiles = append(h.Tiles, tiles...)
	e := Evaluator{Queue: EvaluationQueue{PendingNode: make([]EvaluationNode, 0), Formed: make([]Meld, 0)}}
	formed := e.Evaluate(&h)
	if len(formed) != 5 {
		t.Errorf("should be formed %d", len(formed))
	}
	//B1,B2,B3,B2,B3,B4,B3,B4,B5,B4,B5,B6,B1,B1
}
