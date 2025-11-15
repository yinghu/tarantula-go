package main

import (
	"fmt"
	"slices"

	"gameclustering.com/internal/mj"
)

type MahjongPlayer struct {
	SystemId int64  `json:"SystemId,string"`
	Seat     string `json:"Seat"`
	mj.Hand  `json:"Hand"`
	Auto     bool                  `json:"Auto"`
	B        []mj.Tile             `json:"-"` //bamboo
	C        []mj.Tile             `json:"-"` //character
	D        []mj.Tile             `json:"-"` //dots
	HE       []mj.Tile             `json:"-"` //east
	HS       []mj.Tile             `json:"-"` //south
	HW       []mj.Tile             `json:"-"` //west
	HN       []mj.Tile             `json:"-"` //north
	R        []mj.Tile             `json:"-"` //red
	G        []mj.Tile             `json:"-"` //green
	W        []mj.Tile             `json:"-"` //white
}


func (mp *MahjongPlayer) OnDraw(t mj.Tile) {
	if !mp.Auto {
		return
	}
	switch t.Suit {
	case mj.BAMBOO:
		mp.B = append(mp.B, t)
	case mj.CHARACTER:
		mp.C = append(mp.C, t)
	case mj.DOTS:
		mp.D = append(mp.D, t)
	case mj.HORNOR:
		switch t.Name() {
		case mj.EAST:
			mp.HE = append(mp.HE, t)
		case mj.SOUTH:
			mp.HS = append(mp.HS, t)
		case mj.WEST:
			mp.HW = append(mp.HW, t)
		case mj.NORTH:
			mp.HN = append(mp.HN, t)
		case mj.RED:
			mp.R = append(mp.R, t)
		case mj.GREEN:
			mp.G = append(mp.G, t)
		case mj.WHITE:
			mp.W = append(mp.W, t)
		}
	}
}
func (mp *MahjongPlayer) OnDrop(t mj.Tile) {
	if !mp.Auto {
		return
	}
	switch t.Suit {
	case mj.BAMBOO:
		for i := range mp.B {
			if mp.B[i] == t {
				mp.B = slices.Delete(mp.B, i, i)
			}
		}
	case mj.CHARACTER:
		for i := range mp.C {
			if mp.C[i] == t {
				mp.C = slices.Delete(mp.C, i, i)
			}
		}
	case mj.DOTS:
		for i := range mp.D {
			if mp.D[i] == t {
				mp.D = slices.Delete(mp.D, i, i)
			}
		}
	case mj.HORNOR:
		switch t.Name() {
		case mj.EAST:
			for i := range mp.HE {
				if mp.HE[i] == t {
					mp.HE = slices.Delete(mp.HE, i, i)
				}
			}
		case mj.SOUTH:
			for i := range mp.HS {
				if mp.HS[i] == t {
					mp.HS = slices.Delete(mp.HS, i, i)
				}
			}
		case mj.WEST:
			for i := range mp.HW {
				if mp.HW[i] == t {
					mp.HW = slices.Delete(mp.HW, i, i)
				}
			}
		case mj.NORTH:
			for i := range mp.HN {
				if mp.HN[i] == t {
					mp.HN = slices.Delete(mp.HN, i, i)
				}
			}
		case mj.RED:
			for i := range mp.R {
				if mp.R[i] == t {
					mp.R = slices.Delete(mp.R, i, i)
				}
			}
		case mj.GREEN:
			for i := range mp.G {
				if mp.G[i] == t {
					mp.G = slices.Delete(mp.G, i, i)
				}
			}
		case mj.WHITE:
			for i := range mp.W {
				if mp.W[i] == t {
					mp.W = slices.Delete(mp.W, i, i)
				}
			}
		}
	}
}
func (mp *MahjongPlayer) OnKnog(t mj.Tile) {
	if !mp.Auto {
		return
	}
	switch t.Suit {
	case mj.BAMBOO:
		mp.B = append(mp.B, t)
	case mj.CHARACTER:
		mp.C = append(mp.C, t)
	case mj.DOTS:
		mp.D = append(mp.D, t)
	case mj.HORNOR:
		switch t.Name() {
		case mj.EAST:
			mp.HE = append(mp.HE, t)
		case mj.SOUTH:
			mp.HS = append(mp.HS, t)
		case mj.WEST:
			mp.HW = append(mp.HW, t)
		case mj.NORTH:
			mp.HN = append(mp.HN, t)
		case mj.RED:
			mp.R = append(mp.R, t)
		case mj.GREEN:
			mp.G = append(mp.G, t)
		case mj.WHITE:
			mp.W = append(mp.W, t)
		}
	}
}
func (mp *MahjongPlayer) OnFormed(m mj.Meld) {
	fmt.Printf("tile melt %v\n", m)
}
func NewPlayer(seat string) *MahjongPlayer {
	mp := MahjongPlayer{Seat: seat, Auto: true}
	mp.Hand = mj.Hand{Listener: &mp}
	mp.Hand.New()
	mp.B = make([]mj.Tile, 0)
	mp.C = make([]mj.Tile, 0)
	mp.D = make([]mj.Tile, 0)
	mp.HE = make([]mj.Tile, 0)
	mp.HS = make([]mj.Tile, 0)
	mp.HW = make([]mj.Tile, 0)
	mp.HN = make([]mj.Tile, 0)
	mp.R = make([]mj.Tile, 0)
	mp.G = make([]mj.Tile, 0)
	mp.W = make([]mj.Tile, 0)
	return &mp
}
