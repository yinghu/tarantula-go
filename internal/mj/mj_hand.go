package mj

import (
	"fmt"
	"slices"

	"gameclustering.com/internal/core"
)

type HandListener interface {
	OnDraw(t Tile)
	OnDrop(t Tile)
	OnKnog(t Tile)
	OnFormed(m Meld)
}

type Hand struct {
	Formed    []Meld       `json:"Formed"`
	Tiles     []Tile       `json:"Tiles"`
	Flowers   []Tile       `json:"Flowers"`
	MaxClaims int          `json:"-"`
	Listener  HandListener `json:"-"`
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
			if h.Listener == nil {
				return nil
			}
			h.Listener.OnDrop(drop)
			return nil
		}
	}
	return fmt.Errorf("drop not existed %v", drop)
}
func (h *Hand) Discharge(discharged int) error {
	for i := range h.Tiles {
		core.AppLog.Printf("Tile %d\n", h.Tiles[i].Seq)
		if h.Tiles[i].Seq == discharged {
			drop := h.Tiles[i]
			h.Tiles = slices.Delete(h.Tiles, i, i)
			if h.Listener == nil {
				return nil
			}
			h.Listener.OnDrop(drop)
			return nil
		}
	}
	return fmt.Errorf("discharged not existed %d", discharged)
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
	if h.Listener == nil {
		return nil
	}
	h.Listener.OnDraw(t)
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
	if h.Listener == nil {
		return nil
	}
	h.Listener.OnKnog(t)
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
		if h.Listener != nil {
			h.Listener.OnFormed(v)
		}
	}
	return eyeCount == 1 && formed == 5
}

func (h *Hand) TileSize() int {
	return len(h.Tiles)
}

func (h *Hand) Write(buff core.DataBuffer) error {
	sz := len(h.Tiles)
	if err := buff.WriteInt32(int32(sz)); err != nil {
		return err
	}
	for i := range h.Tiles {
		if err := h.Tiles[i].Write(buff); err != nil {
			return err
		}
	}
	return nil
}
