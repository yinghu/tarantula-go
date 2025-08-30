package mahjong

import (
	"slices"
)

type Hand struct {
	Formed    []Meld
	Dots      []Tile
	Bamboo    []Tile
	Character []Tile
	Honor     []Tile
	Flower    []Tile
}

func cmp(a, b Tile) int {
	diff := a.Num - b.Num
	return int(diff)
}

func (h *Hand) New() {
	h.Formed = make([]Meld, 0)
	h.Dots = make([]Tile, 0)
	h.Bamboo = make([]Tile, 0)
	h.Character = make([]Tile, 0)
	h.Honor = make([]Tile, 0)
	h.Flower = make([]Tile, 0)
}

func (h *Hand) Draw(deck *Deck) error {
	t, err := deck.Draw()
	if err != nil {
		return err
	}
	switch t.Suit {
	case BAMBOO:
		h.Bamboo = append(h.Bamboo, t)
		slices.SortFunc(h.Bamboo, cmp)
	case DOTS:
		h.Dots = append(h.Dots, t)
		slices.SortFunc(h.Dots, cmp)
	case CHARACTER:
		h.Character = append(h.Character, t)
		slices.SortFunc(h.Character, cmp)
	case HORNOR:
		h.Honor = append(h.Honor, t)
		slices.SortFunc(h.Honor, cmp)
	default:
		h.Flower = append(h.Flower, t)
	}
	return nil
}

func (h *Hand) Knog(deck *Deck) error {
	t, err := deck.Kong()
	if err != nil {
		return err
	}
	switch t.Suit {
	case BAMBOO:
		h.Bamboo = append(h.Bamboo, t)
		slices.SortFunc(h.Bamboo, cmp)
	case DOTS:
		h.Dots = append(h.Dots, t)
		slices.SortFunc(h.Dots, cmp)
	case CHARACTER:
		h.Character = append(h.Character, t)
		slices.SortFunc(h.Character, cmp)
	case HORNOR:
		h.Honor = append(h.Honor, t)
		slices.SortFunc(h.Honor, cmp)
	default:
		h.Flower = append(h.Flower, t)
	}
	return nil
}
