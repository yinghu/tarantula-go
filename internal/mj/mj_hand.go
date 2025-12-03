package mj

import (
	"fmt"
	"slices"

	"gameclustering.com/internal/core"
)

type HandListener interface {
	OnDraw(t Tile, kong bool)
	OnDiscard(t Tile)
	OnKong(t Tile)
	OnChow(drop Tile, chow Meld)
	OnPung(pung Meld)
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

func (h *Hand) Discard(discarded Tile) error {
	for i := range h.Tiles {
		if h.Tiles[i].Seq == discarded.Seq {

			h.Tiles = slices.Delete(h.Tiles, i, i+1)
			if discarded.Suit == FLOWER {
				h.Flowers = append(h.Flowers, discarded)
			}
			if h.Listener != nil {
				h.Listener.OnDiscard(discarded)
			}
			return nil
		}
	}
	return fmt.Errorf("discarded not existed %v", discarded)
}

func (h *Hand) Chow(drop Tile, c1 Tile, c2 Tile) error {
	c := []int{-1, -1}
	m := make([]Tile, 0)
	for i := range h.Tiles {
		if h.Tiles[i].Seq == c1.Seq && c[0] == -1 {
			c[0] = i
			m = append(m, h.Tiles[i])
			continue
		}
		if h.Tiles[i].Seq == c2.Seq && c[1] == -1 {
			c[1] = i
			m = append(m, h.Tiles[i])
			continue
		}
	}
	if c[0] == -1 || c[1] == -1 {
		return fmt.Errorf("no chow %s %s %s", drop.Name(), c1.Name(), c2.Name())
	}
	m = append(m, drop)
	slices.SortFunc(m, cmp)
	chow := Meld{Tiles: m}
	if !chow.Chow() {
		return fmt.Errorf("not chow %v", drop)
	}
	c[0] = -1
	c[1] = -1
	h.Tiles = slices.DeleteFunc(h.Tiles, func(t Tile) bool {
		if t.Seq == c1.Seq && c[0] == -1 {
			c[0] = t.Seq
			return true
		}
		if t.Seq == c2.Seq && c[1] == -1 {
			c[1] = t.Seq
			return true
		}
		return false
	})
	h.Formed = append(h.Formed, chow)
	if h.Listener != nil {
		h.Listener.OnChow(drop, chow)
		h.Listener.OnFormed(chow)
	}
	return nil
}

func (h *Hand) Pung(drop Tile) error {
	c := 0
	for i := range h.Tiles {
		if h.Tiles[i].Seq == drop.Seq {
			if c == 2 {
				break
			}
			c++
		}
	}
	if c != 2 {
		return fmt.Errorf("no pung %d", c)
	}
	m := Meld{Tiles: []Tile{drop, drop, drop}}
	c = 0
	h.Tiles = slices.DeleteFunc(h.Tiles, func(t Tile) bool {
		if t.Seq == drop.Seq && c < 2 {
			c++
			return true
		}
		return false
	})
	h.Formed = append(h.Formed, m)
	if h.Listener != nil {
		h.Listener.OnPung(m)
		h.Listener.OnFormed(m)
	}
	return nil
}

func (h *Hand) Kong(kong Tile) error {
	mk := Meld{}
	for i := range h.Formed {
		m := h.Formed[i]
		if m.Pung() && m.Tiles[0].Seq == kong.Seq {
			mk.Tiles = m.Tiles
			mk.Concealed = false
			h.Formed = slices.Delete(h.Formed, i, i+1)
			break
		}
	}
	if mk.Pung() { //pung + a draw
		mk.Tiles = append(mk.Tiles, kong)
		h.Formed = append(h.Formed, mk)
		if h.Listener != nil {
			h.Listener.OnKong(kong) //delete
			h.Listener.OnFormed(mk)
		}
		h.Tiles = slices.DeleteFunc(h.Tiles, func(t Tile) bool {
			return t.Seq == kong.Seq
		})
		return nil
	}
	ct := 0
	for i := range h.Tiles {
		t := h.Tiles[i]
		if t.Seq == kong.Seq {
			ct++
		}
	}
	if ct < 3 {
		return fmt.Errorf("no kong %v", kong)
	}
	if ct == 3 {
		mk.Concealed = false
	} else {
		mk.Concealed = true
	}
	mk.Tiles = []Tile{kong, kong, kong, kong}
	h.Formed = append(h.Formed, mk)
	if h.Listener != nil {
		h.Listener.OnKong(kong)
		h.Listener.OnFormed(mk)
	}
	h.Tiles = slices.DeleteFunc(h.Tiles, func(t Tile) bool {
		return t.Seq == kong.Seq
	})
	return nil
}

func (h *Hand) HeadDraw(deck *Deck) error {
	t, err := deck.Draw()
	if err != nil {
		return err
	}
	return h.append(t, false)
}

func (h *Hand) append(t Tile, kong bool) error {
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
	if h.Listener != nil {
		h.Listener.OnDraw(t, kong)
	}
	return nil
}
func (h *Hand) TailDraw(deck *Deck) error {
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
