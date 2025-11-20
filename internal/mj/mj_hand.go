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
	Formed         []Meld       `json:"Formed"`
	Tiles          []Tile       `json:"Tiles"`
	Flowers        []Tile       `json:"Flowers"`
	MaxClaims      int          `json:"-"`
	Listener       HandListener `json:"-"`
	FlowerExcluded bool
	Sorting        bool
}

func (h *Hand) New() {
	h.Formed = make([]Meld, 0)
	h.Tiles = make([]Tile, 0)
	h.Flowers = make([]Tile, 0)
	h.MaxClaims = CLASSIC_MAX_FORMS
}
func (h *Hand) Clear() {
	h.Formed = h.Formed[:0]
	h.Tiles = h.Tiles[:0]
	h.Flowers = h.Flowers[:0]
}

func (h *Hand) Drop(drop Tile) error {
	for i := range h.Tiles {
		if h.Tiles[i] == drop {
			h.Tiles = slices.Delete(h.Tiles, i, i+1)
			if drop.Suit == FLOWER {
				h.Flowers = append(h.Flowers, drop)
			}
			if h.Listener == nil {
				return nil
			}
			h.Listener.OnDrop(drop)
			return nil
		}
	}
	return fmt.Errorf("drop not existed %v", drop)
}
func (h *Hand) Discharge(discharged int) (Tile, error) {
	for i := range h.Tiles {
		if h.Tiles[i].Seq == discharged {
			drop := h.Tiles[i]
			h.Tiles = slices.Delete(h.Tiles, i, i+1)
			if drop.Suit == FLOWER {
				h.Flowers = append(h.Flowers, drop)
			}
			if h.Listener == nil {
				return drop, nil
			}
			h.Listener.OnDrop(drop)
			return drop, nil
		}
	}
	return Tile{}, fmt.Errorf("discharged not existed %d", discharged)
}

func (h *Hand) Chow() error {
	return nil
}

func (h *Hand) Pung() error {
	return nil
}

func (h *Hand) Draw(deck *Deck) error {
	t, err := deck.Draw()
	if err != nil {
		return err
	}
	return h.append(t, false)
}

func (h *Hand) append(t Tile, knog bool) error {
	if h.FlowerExcluded {
		switch t.Suit {
		case FLOWER:
			h.Flowers = append(h.Flowers, t)
		default:
			h.Tiles = append(h.Tiles, t)
			if h.Sorting {
				slices.SortFunc(h.Tiles, cmp)
			}
		}
	} else {
		h.Tiles = append(h.Tiles, t)
		if h.Sorting {
			slices.SortFunc(h.Tiles, cmp)
		}
	}
	if h.Listener == nil {
		return nil
	}
	if knog {
		h.Listener.OnKnog(t)
	} else {
		h.Listener.OnDraw(t)
	}
	return nil
}
func (h *Hand) Knog(deck *Deck) error {
	t, err := deck.Kong()
	if err != nil {
		return err
	}
	return h.append(t, true)
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
	sz = len(h.Flowers)
	if err := buff.WriteInt32(int32(sz)); err != nil {
		return err
	}
	for i := range h.Flowers {
		if err := h.Flowers[i].Write(buff); err != nil {
			return err
		}
	}
	sz = len(h.Formed)
	if err := buff.WriteInt32(int32(sz)); err != nil {
		return err
	}
	for i := range h.Formed {
		if err := h.Formed[i].Write(buff); err != nil {
			return err
		}
	}
	return nil
}
