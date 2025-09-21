package mahjong

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
	Hand  []Tile
	Index map[int]TileIndex
}

func (h *HandIndex) From(tiles []Tile) {
	h.Hand = tiles
	slices.SortFunc(h.Hand, cmp)
	h.Index = make(map[int]TileIndex)
	for _, v := range tiles {
		s, exists := h.Index[v.Seq]
		if exists {
			s.Count++
			h.Index[v.Seq] = s
		} else {
			h.Index[v.Seq] = TileIndex{Count: 1, Used: 0, Suit: v}
		}
	}
}

func (h *HandIndex) Kong() []Meld {
	h.reset()
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

func (h *HandIndex) Pong() []Meld {
	h.reset()
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
	h.reset()
	nodes := make([]Meld, 0)
	for s, c := range h.Index {
		nc, exsits := h.Index[s+1]
		if !exsits || nc.Count-nc.Used == 0 {
			continue
		}
		nb, exsits := h.Index[s+2]
		if !exsits || nc.Count-nc.Used == 0 {
			continue
		}
		c.Used++
		h.Index[s] = c
		nc.Used++
		h.Index[s+1] = nc
		nb.Used++
		h.Index[s+2] = nb
		tiles := []Tile{c.Suit, nc.Suit, nb.Suit}
		nodes = append(nodes, Meld{Tiles: tiles})
	}
	return nodes
}

func (h *HandIndex) Eye() (Meld, error) {
	h.reset()
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
	h.reset()
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

func (h *HandIndex) reset() {
	for s, c := range h.Index {
		c.Used = 0
		h.Index[s] = c
	}
}
