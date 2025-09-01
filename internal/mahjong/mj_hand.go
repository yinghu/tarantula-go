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
	Next    TileSet
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
	h.Next = h.NewTileSet(THREE_SET)
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
	err := h.evaluate()
	if err != nil {
		fmt.Printf("no match %s\n", err.Error())
		return false
	}
	var eyeCount int
	var formed int
	for _, v := range h.Formed {
		fmt.Printf("X Set size : %v\n", v.Tiles)
		if v.Eye() {
			eyeCount++
		}
		formed++
	}
	return eyeCount == 1 && formed == 5
}

func (h *Hand) evaluate() error {
	for h.TileSize() > 0 {
		next := h.Next
		for {
			if next.Full() {
				h.Formed = append(h.Formed, next.Formed())
				h.Next = h.NewTileSet(THREE_SET)
				break
			}
			t := h.PopTile()
			if next.Allowed(t) {
				next.Append(t)
				break
			}
			fallBack, err := next.Next(h, t)
			for err == nil {
				if fallBack.Allowed(t) {
					fallBack.Append(t)
					h.Next = fallBack
					break
				}
				fallBack, err = fallBack.Next(h, t)
			}
			break
		}
	}
	//last formed
	h.Formed = append(h.Formed, h.Next.Formed())
	return nil
}

func (h *Hand) NewTileSet(id int) TileSet {
	var tset TileSet
	switch id {
	case FOUR_SET:
		tset = NewFourTileSet()
	case THREE_SET:
		tset = NewThreeTileSet()
	case SEQ_SET:
		tset = NewSequenceTileSet()
	case TWO_SET:
		tset = NewTwoTileSet()
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

func (h *Hand) TileSize() int {
	return len(h.Tiles)
}
