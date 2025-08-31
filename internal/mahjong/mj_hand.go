package mahjong

import (
	"slices"
	"strings"
)

type Hand struct {
	Formed  []Meld
	Tiles   []Tile
	Flowers []Tile
	Pending []TileSet
	Sn      int
}

func cmp(a, b Tile) int {
	if a.Suit == b.Suit {
		diff := a.Rank - b.Rank
		return int(diff)
	}
	return strings.Compare(a.Suit, b.Suit)
}

func (h *Hand) New() {
	h.Formed = make([]Meld, 0)
	h.Tiles = make([]Tile, 0)
	h.Flowers = make([]Tile, 0)
	h.Pending = make([]TileSet, 0)
	h.Sn = 1
}

func (h *Hand) Draw(deck *Deck) error {
	t, err := deck.Draw()
	if err != nil {
		return err
	}
	switch t.Suit {
	case FLOWER:
		h.Flowers = append(h.Flowers, t)
	default:
		h.Tiles = append(h.Tiles, t)
		slices.SortFunc(h.Tiles, cmp)
	}
	return nil
}

func (h *Hand) Knog(deck *Deck) error {
	t, err := deck.Kong()
	if err != nil {
		return err
	}
	switch t.Suit {
	case FLOWER:
		h.Flowers = append(h.Flowers, t)
	default:
		h.Tiles = append(h.Tiles, t)
		slices.SortFunc(h.Tiles, cmp)
	}
	return nil
}

func (h *Hand) Mahjong() bool {
	h.Pending = append(h.Pending, NewFourTileSet(h.Sn))
	h.Sn++
	h.evaluate()
	return false
}

func (h *Hand) evaluate() {
	var t Tile
	for {
		sz := len(h.Tiles)
		if sz == 0 {
			return
		}
		t = h.Tiles[0]
		h.Tiles = h.Tiles[1:]
		tset := h.Pending[0]
		if tset.Full() {
			if len(h.Tiles) > 2 {
				h.Pending = slices.Insert(h.Pending, 0, NewFourTileSet(h.Sn))
				h.Sn++
			} else {
				h.Pending = slices.Insert(h.Pending, 0, NewThreeTileSet(h.Sn))
				h.Sn++
			}
			h.Pending[0].Append(t)
		} else if tset.Allowed(t) {
			tset.Append(t)
		} else {
			h.Tiles = slices.Insert(h.Tiles, 0, t)
			if !h.redo() {
				return
			}
		}
	}
}
func (h *Hand) redo() bool {
	if len(h.Pending) == 0 {
		return false
	}
	tset1 := h.Pending[0]
	h.Pending = h.Pending[1:]
	tset2 := tset1.Fallback(h)
	if tset1.Sequence() != tset2.Sequence() {
		h.Pending = slices.Insert(h.Pending, 0, tset2)
	} else {
		for {
			if tset1.Sequence() == tset2.Sequence() && len(h.Pending) != 0 {
				tset1 = h.Pending[0]
				h.Pending = h.Pending[1:]
				tset2 = tset1.Fallback(h)
			} else {
				break
			}
		}
		if tset1.Sequence() == tset2.Sequence() {
			return false
		} else {
			h.Pending = slices.Insert(h.Pending, 0, tset2)
		}

	}
	return true
}
