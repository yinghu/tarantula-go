package mj

import (
	"fmt"
	"slices"
)

func cmp(a, b Tile) int {
	return a.Seq - b.Seq
}

type TileIndex struct {
	Suit  Tile
	Count int
	Used  int
}


type HandIndex struct {
	Hand     []Tile
	Index    map[int]TileIndex
	tx       []int
}

func (h *HandIndex) From(tiles []Tile) {
	h.Hand = tiles
	h.Index = make(map[int]TileIndex)
	h.tx = make([]int, 0)
	for i := range tiles {
		v := tiles[i]
		s, exists := h.Index[v.Seq]
		if exists {
			s.Count++
			h.Index[v.Seq] = s
		} else {
			h.Index[v.Seq] = TileIndex{Count: 1, Used: 0, Suit: v}
			h.tx = append(h.tx, v.Seq)
		}
	}
	slices.Sort(h.tx)
}

func (h *HandIndex) Kong() []Meld {
	h.Reset()
	nodes := make([]Meld, 0)
	for s, c := range h.Index {
		if c.Count-c.Used == 4 {
			tiles := []Tile{c.Suit, c.Suit, c.Suit, c.Suit}
			nodes = append(nodes, Meld{Tiles: tiles})
			c.Used = 4
			h.Index[s] = c
		}
	}
	return nodes
}

func (h *HandIndex) Pung() []Meld {
	h.Reset()
	nodes := make([]Meld, 0)
	for s, c := range h.Index {
		if c.Count-c.Used >= 3 {
			tiles := []Tile{c.Suit, c.Suit, c.Suit}
			nodes = append(nodes, Meld{Tiles: tiles})
			c.Used += 3
			h.Index[s] = c
		}
	}
	return nodes
}

func (h *HandIndex) Chow() []Meld {
	h.Reset()
	nodes := make([]Meld, 0)
	for i := range h.tx {
		s := h.tx[i]
		c, exists := h.Index[s]
		if !exists {
			continue
		}
		nc, exsits := h.Index[s+1]
		if !exsits {
			continue
		}
		nb, exsits := h.Index[s+2]
		if !exsits {
			continue
		}
		tiles := []Tile{c.Suit, nc.Suit, nb.Suit}
		nodes = append(nodes, Meld{Tiles: tiles})
	}
	return nodes
}

func (h *HandIndex) Eye() (Meld, error) {
	h.Reset()
	m := Meld{}
	for s, c := range h.Index {
		if c.Count-c.Used >= 2 {
			tiles := []Tile{c.Suit, c.Suit}
			c.Used += 2
			h.Index[s] = c
			m.Tiles = tiles
			return m, nil
		}
	}
	return m, fmt.Errorf("no eye")
}



func (h *HandIndex) AfterFormed(m Meld) []Tile {
	h.Reset()
	remaining := make([]Tile, 0)
	for _, t := range m.Tiles {
		c := h.Index[t.Seq]
		c.Used++
		h.Index[t.Seq] = c
	}
	for _, c := range h.Index {
		ct := c.Used
		for ct < c.Count {
			remaining = append(remaining, c.Suit)
			ct++
		}
	}
	return remaining
}

func (h *HandIndex) Reset() {
	for s, c := range h.Index {
		c.Used = 0
		h.Index[s] = c
	}
}
/**
func (h *HandIndex) Eyes() []Meld {
	h.reset()
	mk := make([]Meld, 0)
	for s, c := range h.Index {
		if c.Count-c.Used >= 2 {
			tiles := []Tile{c.Suit, c.Suit}
			c.Used += 2
			h.Index[s] = c
			mk = append(mk, Meld{Tiles: tiles})

		}
	}
	return mk
}

func (h *HandIndex) CheckChow(c Tile) []Meld {
	mk := make([]Meld, 0)
	r1, r1e := h.Index[c.Seq+1]
	r2, r2e := h.Index[c.Seq+2]
	if r1e && r2e {
		r := []Tile{c, r1.Suit, r2.Suit}
		mk = append(mk, Meld{Tiles: r})
	}
	l1, l1e := h.Index[c.Seq-1]
	l2, l2e := h.Index[c.Seq-2]
	if l1e && l2e {
		r := []Tile{l2.Suit, l1.Suit, c}
		mk = append(mk, Meld{Tiles: r})
	}
	m1, m1e := h.Index[c.Seq+1]
	m2, m2e := h.Index[c.Seq-1]
	if m1e && m2e {
		r := []Tile{m2.Suit, c, m1.Suit}
		mk = append(mk, Meld{Tiles: r})
	}
	return mk
}

func (h *HandIndex) CheckPung(c Tile) Meld {
	p, exists := h.Index[c.Seq]
	if !exists || p.Count < 2 {
		return Meld{}
	}
	tl := []Tile{c, c, c}
	return Meld{Tiles: tl}
}

func (h *HandIndex) CheckKong(c Tile) Meld {
	p, exists := h.Index[c.Seq]
	if !exists || p.Count < 3 {
		return Meld{}
	}
	tl := []Tile{c, c, c, c}
	return Meld{Tiles: tl}
}**/
