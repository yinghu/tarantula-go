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
	TN           bool //false draw true discharge
	CT           int  //
}

func (mp *MahjongPlayer) OnError(mt *MahjongTable, err error) {
	if !mp.Auto {
		mr := NewMahjongErrorEvent(mp.SystemId, mt.Id, 100, err.Error())
		mt.Push(&mr)
	}
	mt.Sync <- MahjongPlayToken{Cmd: CMD_TURN_END}
	core.AppLog.Printf("play error %s on %d\n", err.Error(), mp.Seat)
}

func (mp *MahjongPlayer) PlayDice(mt *MahjongTable) {
	oid, _ := mt.Sequence.Id()
	md := NewMahjongTurnEvent(mp.SystemId, oid, CMD_DICE, func() {
		mt.Turn <- MahjongPlayToken{Cmd: CMD_DICE, SystemId: mp.SystemId, Seat: mp.Seat, Id: oid}
	}, func() {
		me := MahjongDiceEvent{Dice1: int32(mt.dice[0]), Dice2: int32(mt.dice[1])}
		mt.Update(&me)
	})
	mt.Update(&md)
	mt.Timer <- &md
}

func (mp *MahjongPlayer) PlayDeal(mt *MahjongTable) {
	oid, _ := mt.Sequence.Id()
	md := NewMahjongTurnEvent(mp.SystemId, oid, CMD_DEAL, func() {
		mt.Turn <- MahjongPlayToken{Cmd: CMD_DEAL, SystemId: mp.SystemId, Seat: mp.Seat, Id: oid}
	}, func() {
		me := NewMahjongHandEvent(mp.SystemId, mp.Hand, mp.PendingKongs)
		mt.Update(&me)
	})
	mt.Update(&md)
	mt.Timer <- &md
}

func (mp *MahjongPlayer) PlayDischarge(mt *MahjongTable, mc MahjongDischargeEvent) {
	core.AppLog.Printf("player discharge %d %v\n", mp.Seat, mp.Auto)
	core.AppLog.Printf("drop %v\n", mc.Opts)
	oid, _ := mt.Sequence.Id()
	md := NewMahjongTurnEvent(mp.SystemId, oid, CMD_SKIP, func() {
		mt.Turn <- MahjongPlayToken{Cmd: CMD_SKIP, SystemId: mp.SystemId, Seat: mp.Seat, Id: oid}
	}, func() {
		mt.Sync <- MahjongPlayToken{Cmd: CMD_TURN_END}
	})
	if !mp.Auto {
		mt.Push(&mc)
		mt.Push(&md)
	}
	mt.Timer <- &md
}


func (mp *MahjongPlayer) Play(mt *MahjongTable) {
	core.AppLog.Printf("player turn %d %v %v\n", mp.Seat, mp.Auto, mp.TN)
	core.AppLog.Printf("player kong %v\n", mp.PendingKongs)
	if len(mp.PendingKongs) > 0 { //knog first
		oid, _ := mt.Sequence.Id()
		k := mp.PendingKongs[0]
		md := NewMahjongTurnEvent(mp.SystemId, oid, CMD_KONG, func() {
			mt.Turn <- MahjongPlayToken{Cmd: CMD_KONG, Seat: mp.Seat, Selected: k, Id: oid}
		}, func() {
			if !mp.Auto {
				me := NewMahjongHandEvent(mp.SystemId, mp.Hand, mp.PendingKongs)
				mt.Push(&me)
			}
			mt.Sync <- MahjongPlayToken{Cmd: CMD_TURN_CONTINUE}
		})
		if !mp.Auto {
			mt.Push(&md)
		}
		mt.Timer <- &md
		return
	}
	if mp.TN {
		claimed := mt.Setup.Mahjong(&mp.Hand)
		if claimed {
			oid, _ := mt.Sequence.Id()
			md := NewMahjongTurnEvent(mp.SystemId, oid, CMD_CLAIM, func() {
				mt.Turn <- MahjongPlayToken{Cmd: CMD_CLAIM, Seat: mp.Seat, Id: oid}
			}, func() {
				//reset to next round
				mc := NewMahjongClaimEvent(mp.SystemId, mt.Id, mp.Seat, claimed, mp.Hand.Formed)
				mt.Push(&mc)
			})
			if !mp.Auto {
				mt.Push(&md)
			}
			return
		}
		oid, _ := mt.Sequence.Id()
		md := NewMahjongTurnEvent(mp.SystemId, oid, CMD_DISCHARGE, func() {
			t := mp.Hand.Tiles[oid%10]
			mt.Turn <- MahjongPlayToken{Cmd: CMD_DISCHARGE, Seat: mp.Seat, Selected: t.Seq, Id: oid}
		}, func() {
			mp.TN = false
			if !mp.Auto {
				me := NewMahjongHandEvent(mp.SystemId, mp.Hand, mp.PendingKongs)
				mt.Push(&me)
			}
			mt.Sync <- MahjongPlayToken{Cmd: CMD_TURN_END}
		})
		if !mp.Auto {
			mt.Push(&md)
		}
		mt.Timer <- &md
	} else {
		oid, _ := mt.Sequence.Id()
		md := NewMahjongTurnEvent(mp.SystemId, oid, CMD_DRAW, func() {
			mt.Turn <- MahjongPlayToken{Cmd: CMD_DRAW, Seat: mp.Seat, Id: oid}
		}, func() {
			mp.TN = true
			if !mp.Auto {
				me := NewMahjongHandEvent(mp.SystemId, mp.Hand, mp.PendingKongs)
				mt.Push(&me)
			}
			mt.Sync <- MahjongPlayToken{Cmd: CMD_TURN_CONTINUE}
		})
		if !mp.Auto {
			mt.Push(&md)
		}
		mt.Timer <- &md
	}
}

func (mp *MahjongPlayer) CheckDischarge(seat int, drop mj.Tile, chow bool) []mj.Meld {
	switch drop.Suit {
	case mj.BAMBOO:
		return mp.checkMatch(mp.B, drop, chow)
	case mj.CHARACTER:
		return mp.checkMatch(mp.C, drop, chow)
	case mj.DOTS:
		return mp.checkMatch(mp.D, drop, chow)
	case mj.HORNOR:
		switch drop.Name() {
		case mj.EAST:
			return mp.checkMatch(mp.HE, drop, false)
		case mj.SOUTH:
			return mp.checkMatch(mp.HS, drop, false)
		case mj.WEST:
			return mp.checkMatch(mp.HW, drop, false)
		case mj.NORTH:
			return mp.checkMatch(mp.HN, drop, false)
		case mj.RED:
			return mp.checkMatch(mp.R, drop, false)
		case mj.GREEN:
			return mp.checkMatch(mp.G, drop, false)
		case mj.WHITE:
			return mp.checkMatch(mp.W, drop, false)
		}
	}
	return []mj.Meld{}
}

func (mp *MahjongPlayer) checkMatch(seg []mj.Tile, d mj.Tile, c bool) []mj.Meld {
	if len(seg) < 2 {
		return nil
	}
	ix := mj.HandIndex{}
	seg = append(seg, d)
	ix.From(seg)
	m := make([]mj.Meld, 0)
	if c {

		m = append(m, ix.Chow()...)
		m = append(m, ix.Pung()...)
		m = append(m, ix.Kong()...)
		return m
	}
	m = append(m, ix.Pung()...)
	m = append(m, ix.Kong()...)
	return m
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
		mp.PendingKongs = append(mp.PendingKongs, t.Seq)
	}

}
func (mp *MahjongPlayer) OnDischarge(t mj.Tile) {
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
func (mp *MahjongPlayer) OnKong(t mj.Tile) {
	mp.OnDraw(t)
}
func (mp *MahjongPlayer) OnFormed(m mj.Meld) {

}

func (mp *MahjongPlayer) checkKong(clist []mj.Tile, hornor bool) {
	if !hornor {
		ix := mj.HandIndex{}
		ix.From(clist)
		ks := ix.Kong()
		for i := range ks {
			mp.PendingKongs = append(mp.PendingKongs, ks[i].Tiles[0].Seq)
		}
		return
	}
	if len(clist) == 4 {
		mp.PendingKongs = append(mp.PendingKongs, clist[0].Seq)
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
	mp.CT = 0
	mp.TN = false
	return &mp
}
