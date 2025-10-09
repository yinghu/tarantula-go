package main

import (
	"fmt"
	"slices"
)

type Hand struct {
	Formed    []Meld
	Tiles     []Tile
	Flowers   []Tile
	MaxClaims int
}

func (h *Hand) New() {
	h.Formed = make([]Meld, 0)
	h.Tiles = make([]Tile, 0)
	h.Flowers = make([]Tile, 0)
	h.MaxClaims = CLASSIC_MAX_FORMS
}

func (h *Hand) Drop(drop Tile) error {
	for i := range h.Tiles {
		if h.Tiles[i] == drop {
			h.Tiles = slices.Delete(h.Tiles, i, i)
			return nil
		}
	}
	return fmt.Errorf("drop not existed %v", drop)
}

func (h *Hand) Draw(deck *Deck) error {
	t, err := deck.Draw()
	if err != nil {
		return err
	}
	return h.draw(t)
}

func (h *Hand) draw(t Tile) error {
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
	e := Evaluator{Queue: EvaluationQueue{PendingNode: make([]EvaluationNode, 0), Formed: make([]Meld, 0)}}
	h.Formed = append(h.Formed, e.Evaluate(h)...)
	var eyeCount int
	var formed int
	for _, v := range h.Formed {
		if v.Eye() {
			eyeCount++
		}
		formed++
	}
	return eyeCount == 1 && formed == 5
}

func (h *Hand) TileSize() int {
	return len(h.Tiles)
}
