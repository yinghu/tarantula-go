package mahjong

import "testing"

func TestHandEval1_1(t *testing.T) {
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
	if h.TileSize() != 14 {
		t.Errorf("hand size should be 14 %d", h.TileSize())
	}
	matched := h.Mahjong()
	if !matched {
		t.Errorf("hand should be a match %v", matched)
	}
}
