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
	Tx       []int
}

func (h *HandIndex) From(tiles []Tile) {
	h.Hand = tiles
	h.Index = make(map[int]TileIndex)
	h.Tx = make([]int, 0)
	for i := range tiles {
		v := tiles[i]
		s, exists := h.Index[v.Seq]
		if exists {
			s.Count++
			h.Index[v.Seq] = s
		} else {
			h.Index[v.Seq] = TileIndex{Count: 1, Used: 0, Suit: v}
			h.Tx = append(h.Tx, v.Seq)
		}
	}
	slices.Sort(h.Tx)
}

func (h *HandIndex) Kong() []Meld {
	h.Reset()
	nodes := make([]Meld, 0)
	for _, c := range h.Index {
		if c.Count-c.Used == 4 {
			tiles := []Tile{c.Suit, c.Suit, c.Suit, c.Suit}
			nodes = append(nodes, Meld{Tiles: tiles})
		}
	}
	return nodes
}

func (h *HandIndex) Pung() []Meld {
	h.Reset()
	nodes := make([]Meld, 0)
	for _, c := range h.Index {
		if c.Count-c.Used >= 3 {
			tiles := []Tile{c.Suit, c.Suit, c.Suit}
			nodes = append(nodes, Meld{Tiles: tiles})
		}
	}
	return nodes
}

func (h *HandIndex) Chow() []Meld {
	h.Reset()
	nodes := make([]Meld, 0)
	for i := range h.Tx {
		s := h.Tx[i]
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
	for _, c := range h.Index {
		if c.Count-c.Used >= 2 {
			tiles := []Tile{c.Suit, c.Suit}
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
