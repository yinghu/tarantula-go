package main

import (
	"slices"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/mj"
)

const (
	TC_C int = 0
	TC_B int = 1
	TC_D int = 2
	TC_H int = 3
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
	TN           bool   //false draw true discharge
	TC           [4]int //
}

func (mp *MahjongPlayer) OnError(mt *MahjongTable, err error) {
	if !mp.Auto {
		mr := NewMahjongErrorEvent(mp.SystemId, mt.Id, 100, err.Error())
		mt.Update(&mr)
	}
	core.AppLog.Printf("play error %s on %d\n", err.Error(), mp.Seat)
}

func (mp *MahjongPlayer) OnSeat(mt *MahjongTable) {
	ms := NewMahjongSitEvent(mp.SystemId, mt.Id, mp.Seat)
	mt.Update(&ms)
}

func (mp *MahjongPlayer) PlayDice(mt *MahjongTable) {
	oid, _ := mt.Sequence.Id()
	md := NewMahjongTurnEvent(mp.SystemId, oid, CMD_DICE, func() {
		mt.Turn <- MahjongPlayToken{Cmd: CMD_DICE, SystemId: mp.SystemId, Seat: mp.Seat, Id: oid}
	}, func(commited bool) {
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
	}, func(commited bool) {
		me := NewMahjongHandEvent(mp.SystemId, mp.Hand, mp.PendingKongs)
		mt.Update(&me)
		mt.Sync <- MahjongPlayToken{Cmd: CMD_TURN_START} //start game from dealer
	})
	mt.Update(&md)
	mt.Timer <- &md
}

func (mp *MahjongPlayer) PlayDiscard(mt *MahjongTable, mc MahjongDiscardEvent) {
	core.AppLog.Printf("player discard seat :%d auto:%v tn:%v\n", mp.Seat, mp.Auto, mp.TN)
	core.AppLog.Printf("drop %v\n", mc.Opts)
	oid, _ := mt.Sequence.Id()
	md := NewMahjongTurnEvent(mp.SystemId, oid, CMD_SKIP, func() {
		mt.Turn <- MahjongPlayToken{Cmd: CMD_SKIP, SystemId: mp.SystemId, Seat: mp.Seat, Id: oid}
	}, func(commited bool) {
		if commited {
			me := NewMahjongHandEvent(mp.SystemId, mp.Hand, mp.PendingKongs)
			mt.Update(&me)
			mp.TN = true
			mt.Sync <- MahjongPlayToken{Cmd: CMD_TURN_PLAYER, Seat: mp.Seat} //full hand turn
		} else {
			mt.Sync <- MahjongPlayToken{Cmd: CMD_TURN_NEXT}
		}
	})
	if !mp.Auto {
		mt.Push(&mc)
		mt.Push(&md)
	}
	mt.Timer <- &md
}

func (mp *MahjongPlayer) Play(mt *MahjongTable) {
	core.AppLog.Printf("player seat: %d auto: %v TN: %v\n", mp.Seat, mp.Auto, mp.TN)
	core.AppLog.Printf("player kong list: %v\n", mp.PendingKongs)
	if len(mp.PendingKongs) > 0 { //knog first
		oid, _ := mt.Sequence.Id()
		k := mp.PendingKongs[0]
		md := NewMahjongTurnEvent(mp.SystemId, oid, CMD_KONG, func() {
			mt.Turn <- MahjongPlayToken{Cmd: CMD_KONG, Seat: mp.Seat, Selected: k, Id: oid}
		}, func(commited bool) {
			if !mp.Auto {
				me := NewMahjongHandEvent(mp.SystemId, mp.Hand, mp.PendingKongs)
				mt.Push(&me)
			}
			mt.Sync <- MahjongPlayToken{Cmd: CMD_TURN_PLAYER, Seat: mp.Seat}
		})
		if !mp.Auto {
			mt.Push(&md)
		}
		mt.Timer <- &md
		return
	}
	if mp.TN { //full hand
		claimed := mt.Setup.Mahjong(&mp.Hand)
		core.AppLog.Printf("player claim seat : %d auto: %v mj: %v\n", mp.Seat, mp.Auto, claimed)
		if claimed {
			oid, _ := mt.Sequence.Id()
			md := NewMahjongTurnEvent(mp.SystemId, oid, CMD_CLAIM, func() {
				mt.Turn <- MahjongPlayToken{Cmd: CMD_CLAIM, Seat: mp.Seat, Id: oid}
			}, func(commited bool) {
				//reset to next round
				if !mp.Auto {
					mc := NewMahjongClaimEvent(mp.SystemId, mt.Id, mp.Seat, claimed, mp.Hand.Formed)
					mt.Push(&mc)
				}
			})
			if !mp.Auto {
				mt.Push(&md)
			}
			mt.Timer <- &md
			return
		}
		oid, _ := mt.Sequence.Id()
		md := NewMahjongTurnEvent(mp.SystemId, oid, CMD_DISCARD, func() {
			t := mp.autoPick()
			mt.Turn <- MahjongPlayToken{Cmd: CMD_DISCARD, Seat: mp.Seat, Selected: t.Seq, Id: oid}
		}, func(commited bool) {
			mp.TN = false
			if !mp.Auto {
				me := NewMahjongHandEvent(mp.SystemId, mp.Hand, mp.PendingKongs)
				mt.Push(&me)
			}
			mt.Sync <- MahjongPlayToken{Cmd: CMD_TURN_NEXT} //next player
		})
		if !mp.Auto {
			mt.Push(&md)
		}
		mt.Timer <- &md
		return
	}
	oid, _ := mt.Sequence.Id()
	md := NewMahjongTurnEvent(mp.SystemId, oid, CMD_DRAW, func() {
		mt.Turn <- MahjongPlayToken{Cmd: CMD_DRAW, Seat: mp.Seat, Id: oid} ///
	}, func(commited bool) {
		mp.TN = true
		if !mp.Auto {
			me := NewMahjongHandEvent(mp.SystemId, mp.Hand, mp.PendingKongs)
			mt.Push(&me)
		}
		mt.Sync <- MahjongPlayToken{Cmd: CMD_TURN_PLAYER, Seat: mp.Seat}
	})
	if !mp.Auto {
		mt.Push(&md)
	}
	mt.Timer <- &md
}

func (mp *MahjongPlayer) CheckDiscard(seat int, drop mj.Tile, chow bool) []mj.Meld {
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
		return []mj.Meld{}
	}
	ix := mj.HandIndex{}
	ix.From(seg)
	m := make([]mj.Meld, 0)
	if c {
		m = append(m, ix.CheckChow(d)...)
	}
	cp := ix.CheckPung(d)
	if cp.Type() == int32(mj.PUNG) {
		m = append(m, cp)
	}
	ck := ix.CheckKong(d)
	if ck.Type() == int32(mj.KNOG) {
		mp.PendingKongs = append(mp.PendingKongs, d.Seq)
		m = append(m, ix.CheckKong(d))
	}
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

func (mp *MahjongPlayer) OnDraw(t mj.Tile, kong bool) {
	switch t.Suit {
	case mj.BAMBOO:
		mp.B = append(mp.B, t)
		mp.checkKong(mp.B, false)
		mp.TC[TC_B]++
	case mj.CHARACTER:
		mp.C = append(mp.C, t)
		mp.checkKong(mp.C, false)
		mp.TC[TC_C]++
	case mj.DOTS:
		mp.D = append(mp.D, t)
		mp.checkKong(mp.D, false)
		mp.TC[TC_D]++
	case mj.HORNOR:
		mp.TC[TC_H]++
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

func (mp *MahjongPlayer) OnDiscard(t mj.Tile) {
	mp.OnDelete(t, false)
}
func (mp *MahjongPlayer) OnDelete(t mj.Tile, kong bool) {
	switch t.Suit {
	case mj.BAMBOO:
		if !kong {
			for i := range mp.B {
				if mp.B[i] == t {
					mp.B = slices.Delete(mp.B, i, i+1)
					mp.TC[TC_B]--
					break
				}
			}
		} else {
			mp.B = slices.DeleteFunc(mp.B, func(d mj.Tile) bool {
				if d.Seq == t.Seq {
					mp.TC[TC_B]--
				}
				return d.Seq == t.Seq
			})
		}

	case mj.CHARACTER:
		if !kong {
			for i := range mp.C {
				if mp.C[i] == t {
					mp.C = slices.Delete(mp.C, i, i+1)
					mp.TC[TC_C]--
					break
				}
			}
		} else {
			mp.C = slices.DeleteFunc(mp.C, func(d mj.Tile) bool {
				if d.Seq == t.Seq {
					mp.TC[TC_C]--
				}
				return d.Seq == t.Seq
			})
		}

	case mj.DOTS:
		if !kong {
			for i := range mp.D {
				if mp.D[i] == t {
					mp.D = slices.Delete(mp.D, i, i+1)
					mp.TC[TC_D]--
					break
				}
			}
		} else {
			mp.D = slices.DeleteFunc(mp.D, func(d mj.Tile) bool {
				if d.Seq == t.Seq {
					mp.TC[TC_D]--
				}
				return d.Seq == t.Seq
			})
		}

	case mj.HORNOR:
		switch t.Name() {
		case mj.EAST:
			if !kong {
				for i := range mp.HE {
					if mp.HE[i] == t {
						mp.HE = slices.Delete(mp.HE, i, i+1)
						mp.TC[TC_H]--
						break
					}
				}
			} else {
				mp.HE = slices.DeleteFunc(mp.HE, func(d mj.Tile) bool {
					if d.Seq == t.Seq {
						mp.TC[TC_H]--
					}
					return d.Seq == t.Seq
				})
			}

		case mj.SOUTH:
			if !kong {
				for i := range mp.HS {
					if mp.HS[i] == t {
						mp.HS = slices.Delete(mp.HS, i, i+1)
						mp.TC[TC_H]--
						break
					}
				}
			} else {
				mp.HS = slices.DeleteFunc(mp.HS, func(d mj.Tile) bool {
					if d.Seq == t.Seq {
						mp.TC[TC_H]--
					}
					return d.Seq == t.Seq
				})
			}
		case mj.WEST:
			if !kong {
				for i := range mp.HW {
					if mp.HW[i] == t {
						mp.HW = slices.Delete(mp.HW, i, i+1)
						mp.TC[TC_H]--
						break
					}
				}
			} else {
				mp.HW = slices.DeleteFunc(mp.HW, func(d mj.Tile) bool {
					if d.Seq == t.Seq {
						mp.TC[TC_H]--
					}
					return d.Seq == t.Seq
				})
			}

		case mj.NORTH:
			if !kong {
				for i := range mp.HN {
					if mp.HN[i] == t {
						mp.HN = slices.Delete(mp.HN, i, i+1)
						mp.TC[TC_H]--
						break
					}
				}
			} else {
				mp.HN = slices.DeleteFunc(mp.HN, func(d mj.Tile) bool {
					if d.Seq == t.Seq {
						mp.TC[TC_H]--
					}
					return d.Seq == t.Seq
				})
			}

		case mj.RED:
			if !kong {
				for i := range mp.R {
					if mp.R[i] == t {
						mp.R = slices.Delete(mp.R, i, i+1)
						mp.TC[TC_H]--
						break
					}
				}
			} else {
				mp.R = slices.DeleteFunc(mp.R, func(d mj.Tile) bool {
					if d.Seq == t.Seq {
						mp.TC[TC_H]--
					}
					return d.Seq == t.Seq
				})
			}

		case mj.GREEN:
			if !kong {
				for i := range mp.G {
					if mp.G[i] == t {
						mp.G = slices.Delete(mp.G, i, i+1)
						mp.TC[TC_H]--
						break
					}
				}
			} else {
				mp.G = slices.DeleteFunc(mp.G, func(d mj.Tile) bool {
					if d.Seq == t.Seq {
						mp.TC[TC_H]--
					}
					return d.Seq == t.Seq
				})
			}

		case mj.WHITE:
			if !kong {
				for i := range mp.W {
					if mp.W[i] == t {
						mp.W = slices.Delete(mp.W, i, i+1)
						mp.TC[TC_H]--
						break
					}
				}
			} else {
				mp.W = slices.DeleteFunc(mp.W, func(d mj.Tile) bool {
					if d.Seq == t.Seq {
						mp.TC[TC_H]--
					}
					return d.Seq == t.Seq
				})
			}
		}

	}
}
func (mp *MahjongPlayer) OnKong(t mj.Tile) {
	mp.OnDelete(t, true)
}
func (mp *MahjongPlayer) OnFormed(m mj.Meld) {
	core.AppLog.Printf("Meld : %s\n", m.Name())
}

func (mp *MahjongPlayer) OnChow(drop mj.Tile, chow mj.Meld) {
	for i := range chow.Tiles {
		if chow.Tiles[i].Seq != drop.Seq {
			mp.OnDiscard(chow.Tiles[i])
		}
	}
}
func (mp *MahjongPlayer) OnPung(pung mj.Meld) {
	mp.OnDiscard(pung.Tiles[0])
	mp.OnDiscard(pung.Tiles[0])
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
	mp.TN = false
	return &mp
}

func (mp *MahjongPlayer) autoPick() mj.Tile {
	ts := CheckDiscard(mp)
	return mj.FromQ(ts)
}
