package main

import (
	"slices"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/mj"
)

type MahjongPlayer struct {
	SystemId     int64 `json:"SystemId,string"`
	Seat         int   `json:"Seat"`
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
	cmd          int
}

func (mp *MahjongPlayer) Setup(cmd int, mt *MahjongTable) {
	oid, _ := mt.Sequence.Id()
	md := NewMahjongTurnEvent(mp.SystemId, oid, cmd, func() {
		mt.Turn <- MahjongPlayToken{Cmd: cmd, SystemId: mp.SystemId, Seat: mp.Seat, Id: oid}
	})
	mt.Push(&md)
	mt.Timer <- &md
}

func (mp *MahjongPlayer) Play(mt *MahjongTable) {
	core.AppLog.Printf("Player Seat %d Auto %v\n", mp.Seat, mp.Auto)
	if !mp.Auto {

		oid, _ := mt.Sequence.Id()
		md := NewMahjongTurnEvent(mp.SystemId, oid, mp.cmd, func() {
			t := mp.Hand.Tiles[0]
			mt.Turn <- MahjongPlayToken{Cmd: mp.cmd, SystemId: mp.SystemId, Seat: mp.Seat, Selected: t.Seq, Id: oid}
		})
		mt.Push(&md)
		mt.Timer <- &md
		return
	}
	t := mp.Hand.Tiles[0]
	mt.Turn <- MahjongPlayToken{Cmd: CMD_DISCHARGE, SystemId: mp.SystemId, Seat: mp.Seat, Selected: t.Seq, Id: 0}
}
func (mp *MahjongPlayer) OnPlayFinished(m *MahjongTable, t MahjongPlayToken, err error) {
	if mp.Auto {
		if err != nil {
			mp.cmd = CMD_DRAW
		}
		core.AppLog.Printf("Auto Hand %d >> %d \n", mp.Seat, mp.Hand.TileSize())
		m.Sync <- MahjongPlayToken{Cmd: CMD_TURN_END}
		return
	}
	if err != nil {
		mp.cmd = CMD_DRAW
		mr := NewMahjongErrorEvent(mp.SystemId, m.Id, 100, err.Error())
		m.Push(&mr)
	} else {
		mt := MahjongHandEvent{H: mp.Hand, K: mp.PendingKongs}
		m.Push(&mt)
		if t.Cmd == CMD_DISCHARGE {
			dz := len(m.Discharged)
			if dz <= 3 {
				md := MahjongDischargeEvent{D: m.Discharged}
				m.Push(&md)
			} else {
				md := MahjongDischargeEvent{D: m.Discharged[(dz - 3):]}
				m.Push(&md)
			}
		}
		m.Sync <- MahjongPlayToken{Cmd: CMD_TURN_END}
	}
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
		if !mp.Auto {
			//core.AppLog.Printf("Knog %d %d\n", mp.SystemId, t.Seq)
			mp.PendingKongs = append(mp.PendingKongs, t.Seq)
			//mt := MahjongKnogEvent{SystemId: mp.SystemId, Knog: mp.PendingKongs}
			//mp.Pusher.Push(&mt)
		}
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
		for i := range mp.B {
			if mp.B[i] == t {
				mp.B = slices.Delete(mp.B, i, i+1)
				break
			}
		}

	case mj.CHARACTER:

		for i := range mp.C {
			if mp.C[i] == t {
				mp.C = slices.Delete(mp.C, i, i+1)
				break
			}
		}

	case mj.DOTS:

		for i := range mp.D {
			if mp.D[i] == t {
				mp.D = slices.Delete(mp.D, i, i+1)
				break
			}
		}

	case mj.HORNOR:
		switch t.Name() {
		case mj.EAST:

			for i := range mp.HE {
				if mp.HE[i] == t {
					mp.HE = slices.Delete(mp.HE, i, i+1)
					break
				}
			}

		case mj.SOUTH:

			for i := range mp.HS {
				if mp.HS[i] == t {
					mp.HS = slices.Delete(mp.HS, i, i+1)
					break
				}
			}

		case mj.WEST:

			for i := range mp.HW {
				if mp.HW[i] == t {
					mp.HW = slices.Delete(mp.HW, i, i+1)
					break
				}
			}

		case mj.NORTH:

			for i := range mp.HN {
				if mp.HN[i] == t {
					mp.HN = slices.Delete(mp.HN, i, i+1)
					break
				}
			}

		case mj.RED:

			for i := range mp.R {
				if mp.R[i] == t {
					mp.R = slices.Delete(mp.R, i, i+1)
					break
				}
			}

		case mj.GREEN:

			for i := range mp.G {
				if mp.G[i] == t {
					mp.G = slices.Delete(mp.G, i, i+1)
					break
				}
			}

		case mj.WHITE:

			for i := range mp.W {
				if mp.W[i] == t {
					mp.W = slices.Delete(mp.W, i, i+1)
					break
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
	//core.AppLog.Printf("Check list %v\n", clist)
	if !hornor {

		return
	}
	if len(clist) == 4 {
		mp.PendingKongs = append(mp.PendingKongs, clist[0].Seq)
		mt := MahjongKnogEvent{Knog: mp.PendingKongs}
		mt.SystemId = mp.SystemId
		mp.Pusher.Push(&mt)
	}

}

func NewPlayer(seat int, sorting bool, pusher event.Pusher) *MahjongPlayer {
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
	mp.cmd = CMD_DISCHARGE
	return &mp
}
