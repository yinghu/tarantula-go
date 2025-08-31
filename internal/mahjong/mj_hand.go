package mahjong

import (
	"fmt"
	"slices"
	"strings"
)

const (
	FOUR_SET  int = 0
	THREE_SET int = 1
	SEQ_SET   int = 2
	TWO_SET   int = 3
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
	slices.SortFunc(h.Tiles, cmp)
	fmt.Printf("%v\n", h.Tiles)
	h.Pending = append(h.Pending, h.NewTileSet(FOUR_SET))
	h.evaluate()
	fmt.Printf("tile set size : %d\n", h.TileSetSize())
	for i := range h.Pending {
		fmt.Printf("X Set size : %d\n", h.Pending[i].Size())
	}
	return false
}

func (h *Hand) evaluate() {
	var t Tile
	for h.TileSize() > 0 {
		t = h.PopTile()
		tset := h.Pending[0]
		if tset.Full() {
			if h.TileSize() > 2 {
				h.PushTileSet(h.NewTileSet(FOUR_SET))
			} else {
				h.PushTileSet(h.NewTileSet(THREE_SET))
			}
			h.Pending[0].Append(t)
		} else if tset.Allowed(t) {
			tset.Append(t)
		} else {
			h.PushTile(t)
			if !h.redo() {
				return
			}
		}
	}
}
func (h *Hand) redo() bool {
	if h.TileSetSize() == 0 {
		return false
	}
	tset1 := h.PopTileSet()
	fmt.Printf("xpop tile set : %d\n", tset1.Sequence())
	tset2 := tset1.Fallback(h)
	if tset1.Sequence() != tset2.Sequence() {
		h.PushTileSet(tset2)
		return true
	}
	for tset1.Sequence() == tset2.Sequence() && h.TileSetSize() != 0 {
		tset1 = h.PopTileSet()
		tset2 = tset1.Fallback(h)
		fmt.Printf("pop tile set : %d\n", tset1.Sequence())
	}
	if tset1.Sequence() == tset2.Sequence() {
		return false
	}
	h.PushTileSet(tset2)
	return true
}

func (h *Hand) NewTileSet(id int) TileSet {
	var tset TileSet
	switch id {
	case FOUR_SET:
		tset = NewFourTileSet(h.Sn)
		h.Sn++
	case THREE_SET:
		tset = NewThreeTileSet(h.Sn)
		h.Sn++
	case SEQ_SET:
		tset = NewSequenceTileSet(h.Sn)
		h.Sn++
	case TWO_SET:
		tset = NewTwoTileSet(h.Sn)
		h.Sn++
	}
	return tset
}

func (h *Hand) PopTile() Tile {
	t := h.Tiles[0]
	h.Tiles = h.Tiles[1:]
	return t
}

func (h *Hand) PushTile(t Tile) {
	h.Tiles = slices.Insert(h.Tiles, 0, t)
}

func (h *Hand) PopTileSet() TileSet {
	t := h.Pending[0]
	h.Pending = h.Pending[1:]
	//fmt.Printf("pop tile set : %d\n", t.Sequence())
	return t
}

func (h *Hand) PushTileSet(t TileSet) {
	//fmt.Printf("push tile set : %d\n", t.Sequence())
	h.Pending = slices.Insert(h.Pending, 0, t)
}

func (h *Hand) TileSize() int {
	return len(h.Tiles)
}

func (h *Hand) TileSetSize() int {
	return len(h.Pending)
}
