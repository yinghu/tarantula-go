package mahjong

import (
	"slices"
	"strings"
)

type Hand struct {
	Formed  []Meld
	Tiles   []Tile
	Flowers []Tile
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
	return false
}
