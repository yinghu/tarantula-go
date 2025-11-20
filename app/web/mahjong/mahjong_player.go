package main

import (
	"slices"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/mj"
)

type MahjongPlayer struct {
	SystemId     int64  `json:"SystemId,string"`
	Seat         string `json:"Seat"`
	mj.Hand      `json:"Hand"`
	Auto         bool      `json:"Auto"`
	B            []mj.Tile `json:"-"` //bamboo
	C            []mj.Tile `json:"-"` //character
	D            []mj.Tile `json:"-"` //dots
	HE           []mj.Tile `json:"-"` //east
	HS           []mj.Tile `json:"-"` //south
	HW           []mj.Tile `json:"-"` //west
	HN           []mj.Tile `json:"-"` //north
	R            []mj.Tile `json:"-"` //red
	G            []mj.Tile `json:"-"` //green
	W            []mj.Tile `json:"-"` //white
	PendingKongs []int
	Pusher       event.Pusher
}

func (mp *MahjongPlayer) Reset() {
	mp.Clear()
	mp.B = mp.B[:0]
	mp.C = mp.C[:0]
	mp.D = mp.D[:0]
	mp.HE = mp.HE[:0]
	mp.HS = mp.HS[:0]
	mp.HW = mp.HW[:0]
	mp.HN = mp.HN[:0]
	mp.R = mp.R[:0]
	mp.G = mp.G[:0]
	mp.W = mp.W[:0]
	mp.PendingKongs = mp.PendingKongs[:0]
}

func (mp *MahjongPlayer) OnDraw(t mj.Tile) {
	switch t.Suit {
	case mj.BAMBOO:
		mp.B = append(mp.B, t)
		mp.checkKong(mp.B, false)
	case mj.CHARACTER:
		mp.C = append(mp.C, t)
		mp.checkKong(mp.C, false)
	case mj.DOTS:
		mp.D = append(mp.D, t)
		mp.checkKong(mp.D, false)
	case mj.HORNOR:
		switch t.Name() {
		case mj.EAST:
			mp.HE = append(mp.HE, t)
			mp.checkKong(mp.HE, true)
		case mj.SOUTH:
			mp.HS = append(mp.HS, t)
			mp.checkKong(mp.HS, true)
		case mj.WEST:
			mp.HW = append(mp.HW, t)
			mp.checkKong(mp.HW, true)
		case mj.NORTH:
			mp.HN = append(mp.HN, t)
			mp.checkKong(mp.HN, true)
		case mj.RED:
			mp.R = append(mp.R, t)
			mp.checkKong(mp.R, true)
		case mj.GREEN:
			mp.G = append(mp.G, t)
			mp.checkKong(mp.G, true)
		case mj.WHITE:
			mp.W = append(mp.W, t)
			mp.checkKong(mp.W, true)
		}
	case mj.FLOWER:
		core.AppLog.Printf("Knog %d %d\n", mp.SystemId, t.Seq)
		mp.PendingKongs = append(mp.PendingKongs, t.Seq)
		mt := MahjongKnogEvent{SystemId: mp.SystemId, Knog: mp.PendingKongs}
		mp.Pusher.Push(&mt)
	}

}
func (mp *MahjongPlayer) OnDrop(t mj.Tile) {
	if !mp.Auto {
		if t.Suit == mj.FLOWER {
			//KNOG EVENT TO PLAYER
			//mp.PendingKongs = append(mp.PendingKongs, t.Seq)
			//mt := MahjongKnogEvent{SystemId: mp.SystemId, Knog: mp.PendingKongs}
			//mp.Pusher.Push(&mt)
		}
	}
	switch t.Suit {
	case mj.BAMBOO:
		sz := len(mp.B)
		if sz > 0 {
			for i := range mp.B {
				if mp.B[i] == t {
					if i == sz-1 {
						mp.B = mp.B[:sz-1]
					} else {
						mp.B = slices.Delete(mp.B, i, i+1)
					}
					break
				}
			}
		}
	case mj.CHARACTER:
		sz := len(mp.C)
		if sz > 0 {
			for i := range mp.C {
				if mp.C[i] == t {
					if i == sz-1 {
						mp.C = mp.C[:sz-1]
					} else {
						mp.C = slices.Delete(mp.C, i, i+1)
					}
					break
				}
			}
		}
	case mj.DOTS:
		sz := len(mp.D)
		if sz > 0 {
			for i := range mp.D {
				if mp.D[i] == t {
					if i == sz-1 {
						mp.D = mp.D[:sz-1]
					} else {
						mp.D = slices.Delete(mp.D, i, i+1)
					}
					break
				}
			}
		}
	case mj.HORNOR:
		switch t.Name() {
		case mj.EAST:
			sz := len(mp.HE)
			if sz > 0 {
				for i := range mp.HE {
					if mp.HE[i] == t {
						if i == sz-1 {
							mp.HE = mp.HE[:sz-1]
						} else {
							mp.HE = slices.Delete(mp.HE, i, i+1)
						}
						break
					}
				}
			}
		case mj.SOUTH:
			sz := len(mp.HS)
			if sz > 0 {
				for i := range mp.HS {
					if mp.HS[i] == t {
						if i == sz-1 {
							mp.HS = mp.HS[:sz-1]
						} else {
							mp.HS = slices.Delete(mp.HS, i, i+1)
						}
						break
					}
				}
			}
		case mj.WEST:
			sz := len(mp.HW)
			if sz > 0 {
				for i := range mp.HW {
					if mp.HW[i] == t {
						if i == sz-1 {
							mp.HW = mp.HW[:sz-1]
						} else {
							mp.HW = slices.Delete(mp.HW, i, i+1)
						}
						break
					}
				}
			}
		case mj.NORTH:
			sz := len(mp.HN)
			if sz > 0 {
				for i := range mp.HN {
					if mp.HN[i] == t {
						if i == sz-1 {
							mp.HN = mp.HN[:sz-1]
						} else {
							mp.HN = slices.Delete(mp.HN, i, i+1)
						}
						break
					}
				}
			}
		case mj.RED:
			sz := len(mp.R)
			if sz > 0 {
				for i := range mp.R {
					if mp.R[i] == t {
						if i == sz-1 {
							mp.R = mp.R[:sz-1]
						} else {
							mp.R = slices.Delete(mp.R, i, i+1)
						}
						break
					}
				}
			}
		case mj.GREEN:
			sz := len(mp.G)
			if sz > 0 {
				for i := range mp.G {
					if mp.G[i] == t {
						if i == sz-1 {
							mp.G = mp.G[:sz-1]
						} else {
							mp.G = slices.Delete(mp.G, i, i+1)
						}
						break
					}
				}
			}
		case mj.WHITE:
			sz := len(mp.W)
			if sz > 0 {
				for i := range mp.W {
					if mp.W[i] == t {
						if i == sz-1 {
							mp.W = mp.W[:sz-1]
						} else {
							mp.W = slices.Delete(mp.W, i, i+1)
						}
						break
					}
				}
			}
		}
	}
}
func (mp *MahjongPlayer) OnKnog(t mj.Tile) {
	mp.OnDraw(t)
}
func (mp *MahjongPlayer) OnFormed(m mj.Meld) {

}

func (mp *MahjongPlayer) checkKong(clist []mj.Tile, hornor bool) {
	core.AppLog.Printf("Check list %v\n", clist)
	if !hornor {

		return
	}
	if len(clist) == 4 {
		mp.PendingKongs = append(mp.PendingKongs, clist[0].Seq)
		mt := MahjongKnogEvent{SystemId: mp.SystemId, Knog: mp.PendingKongs}
		mp.Pusher.Push(&mt)
	}

}

func NewPlayer(seat string, sorting bool, pusher event.Pusher) *MahjongPlayer {
	mp := MahjongPlayer{Seat: seat, Auto: true, Pusher: pusher}
	mp.Hand = mj.Hand{Listener: &mp, Sorting: sorting}
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
	mp.PendingKongs = make([]int, 0)
	return &mp
}
