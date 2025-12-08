package main

import (
	"slices"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/mj"
)

func hmp(a, b HandSegmenet) int {
	return a.C - b.C
}

type HandSegmenet struct {
	T int
	C int
	mj.HandIndex
	Used map[int]int
}

func (s *HandSegmenet) OnIndex(t mj.TileIndex) {
	//core.AppLog.Printf("Tile usage %s : CT : %d USED : %d\n", t.Suit.Name(), t.Count, t.Used)
	s.Used[t.Suit.Seq] = t.Used
}

func (h *HandSegmenet) CheckChow(c mj.Tile) []mj.Meld {
	mk := make([]mj.Meld, 0)
	r1, r1e := h.Index[c.Seq+1]
	r2, r2e := h.Index[c.Seq+2]
	if r1e && r2e {
		r := []mj.Tile{c, r1.Suit, r2.Suit}
		mk = append(mk, mj.Meld{Tiles: r})
	}
	l1, l1e := h.Index[c.Seq-1]
	l2, l2e := h.Index[c.Seq-2]
	if l1e && l2e {
		r := []mj.Tile{l2.Suit, l1.Suit, c}
		mk = append(mk, mj.Meld{Tiles: r})
	}
	m1, m1e := h.Index[c.Seq+1]
	m2, m2e := h.Index[c.Seq-1]
	if m1e && m2e {
		r := []mj.Tile{m2.Suit, c, m1.Suit}
		mk = append(mk, mj.Meld{Tiles: r})
	}
	return mk
}

func (h *HandSegmenet) CheckPung(c mj.Tile) mj.Meld {
	p, exists := h.Index[c.Seq]
	if !exists || p.Count < 2 {
		return mj.Meld{}
	}
	tl := []mj.Tile{c, c, c}
	return mj.Meld{Tiles: tl}
}

func (h *HandSegmenet) CheckKong(c mj.Tile) mj.Meld {
	p, exists := h.Index[c.Seq]
	if !exists || p.Count < 3 {
		return mj.Meld{}
	}
	tl := []mj.Tile{c, c, c, c}
	return mj.Meld{Tiles: tl}
}

func (h *HandSegmenet) CheckDiscard(mp *MahjongPlayer) int {
	if mp.TC[TC_H] > 0 { //hornor first
		if len(mp.HE) == 1 {
			return mp.HE[0].Seq
		}
		if len(mp.HS) == 1 {
			return mp.HS[0].Seq
		}
		if len(mp.HW) == 1 {
			return mp.HW[0].Seq
		}
		if len(mp.HN) == 1 {
			return mp.HN[0].Seq
		}
		if len(mp.R) == 1 {
			return mp.R[0].Seq
		}
		if len(mp.G) == 1 {
			return mp.G[0].Seq
		}
		if len(mp.W) == 1 {
			return mp.W[0].Seq
		}
	}
	tcs := make([]HandSegmenet, 0)
	tcs = append(tcs, HandSegmenet{T: TC_C, C: mp.TC[TC_C]})
	tcs = append(tcs, HandSegmenet{T: TC_B, C: mp.TC[TC_B]})
	tcs = append(tcs, HandSegmenet{T: TC_D, C: mp.TC[TC_D]})
	slices.SortFunc(tcs, hmp)
	for i := range tcs {
		c := tcs[i]
		if c.C == 0 {
			continue
		}
		core.AppLog.Printf("Discard from suit %d count %d\n", c.T, c.C)
		if c.T == TC_C {
			dis := h.discard(mp.C)
			if dis > 0 {
				return dis
			}
		}
		if c.T == TC_B {
			dis := h.discard(mp.B)
			if dis > 0 {
				return dis
			}
		}
		if c.T == TC_D {
			dis := h.discard(mp.D)
			if dis > 0 {
				return dis
			}
		}
	}
	core.AppLog.Printf("oops something wrong")
	return mp.Hand.Tiles[0].Seq
}

func (h *HandSegmenet) discard(seg []mj.Tile) int {
	core.AppLog.Printf("Segement %v\n", seg)
	h.From(seg)
	h.Kong()
	h.Pung()
	h.Chow()
	for s := range h.Index {
		if h.Used[s] == 0 {
			return s
		}
	}
	return seg[0].Seq
}
