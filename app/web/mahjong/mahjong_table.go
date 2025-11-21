package main

import (
	"fmt"
	"slices"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/mj"
)

const (
	SEAT_E int = 0
	SEAT_S int = 1
	SEAT_W int = 2
	SEAT_N int = 3

	HAND_SIZE_THRESHHOLD int = 13
	TURN_TICKER_SECONDS  int = 30

	SOLO int = 1
)

type MahjongTable struct {
	Id           int64             `json:"Id,string"`
	Setup        mj.ClassicMahjong `json:"-"`
	Players      [4]*MahjongPlayer `json:"Players"`
	Pts          int               `json:"Pts"`
	Discharged   []mj.Tile         `json:"Discharged"`
	Started      bool
	Solo         bool
	Turn         chan MahjongPlayToken `json:"-"`
	event.Pusher `json:"-"`
}

func (m *MahjongTable) New() {
	m.Setup.New()
	m.Players[SEAT_E] = NewPlayer(SEAT_E, true, m)
	m.Players[SEAT_S] = NewPlayer(SEAT_S, true, m)
	m.Players[SEAT_W] = NewPlayer(SEAT_W, true, m)
	m.Players[SEAT_N] = NewPlayer(SEAT_N, true, m)
	m.Discharged = make([]mj.Tile, 0)
	m.Turn = make(chan MahjongPlayToken, 3)
}

func (m *MahjongTable) Reset() {
	m.Setup.New()
	m.Players[SEAT_E].Reset()
	m.Players[SEAT_S].Reset()
	m.Players[SEAT_W].Reset()
	m.Players[SEAT_N].Reset()
	m.Discharged = m.Discharged[:0]
}

func (m *MahjongTable) Play() {
	tick := time.NewTicker(time.Duration(TURN_TICKER_SECONDS) * time.Second)
	running := true
	nextTurn := SEAT_E
	for running {
		select {
		case <-tick.C:
			core.AppLog.Printf("Turn %d\n", nextTurn%4)
			m.Players[nextTurn%4].StartTurn(m)
		case t := <-m.Turn:
			if t.Cmd == CMD_END {
				running = false
				break
			}
			if t.Cmd == CMD_END_TURN {
				nextTurn++
				continue
			}
			core.AppLog.Printf("Token seat: %d selected : %d cmd: %d \n", t.Seat, t.Selected, t.Cmd)
			switch t.Cmd {
			case CMD_SIT:
				err := m.Sit(t.SystemId, t.Seat)
				if err != nil {
					mr := MahjongErrorEvent{SystemId: t.SystemId, TableId: m.Id, Code: 100, Message: err.Error()}
					m.Push(&mr)
				} else {
					mt := MahjongSitEvent{SystemId: t.SystemId, TableId: m.Id, Seat: int32(t.Seat)}
					m.Push(&mt)
				}
			case CMD_DICE:
				dice := m.Dice()
				mt := MahjongDiceEvent{Dice1: int32(dice[0]), Dice2: int32(dice[1])}
				m.Push(&mt)
			case CMD_DEAL:
				m.Deal()
				mt := MahjongHandEvent{H: m.Players[t.Seat].Hand, K: m.Players[t.Seat].PendingKongs}
				m.Push(&mt)
			case CMD_DRAW:
				err := m.Draw(t.Seat)
				if err != nil {
					mr := MahjongErrorEvent{SystemId: t.SystemId, TableId: m.Id, Code: 100, Message: err.Error()}
					m.Push(&mr)
				} else {
					mt := MahjongHandEvent{H: m.Players[t.Seat].Hand, K: m.Players[t.Seat].PendingKongs}
					m.Push(&mt)
				}
			case CMD_KONG:
				err := m.Knog(t.Seat, t.Selected)
				if err != nil {
					mr := MahjongErrorEvent{SystemId: t.SystemId, TableId: m.Id, Code: 100, Message: err.Error()}
					m.Push(&mr)
				} else {
					mt := MahjongHandEvent{H: m.Players[t.Seat].Hand, K: m.Players[t.Seat].PendingKongs}
					m.Push(&mt)
				}
			case CMD_DISCHARGE:
				err := m.Discharge(t.Seat, t.Selected)
				if err != nil {
					mr := MahjongErrorEvent{SystemId: t.SystemId, TableId: m.Id, Code: 100, Message: err.Error()}
					m.Push(&mr)
				} else {
					mt := MahjongHandEvent{H: m.Players[t.Seat].Hand, K: m.Players[t.Seat].PendingKongs}
					m.Push(&mt)
					dz := len(m.Discharged)
					if dz <= 3 {
						md := MahjongDischargeEvent{D: m.Discharged}
						m.Push(&md)
					} else {
						md := MahjongDischargeEvent{D: m.Discharged[(dz - 3):]}
						m.Push(&md)
					}
				}
			case CMD_CLAIM:
				claimed := m.Claim(t.Seat)
				mt := MahjongClaimEvent{SystemId: t.SystemId, TableId: m.Id, Seat: int32(t.Seat), Claimed: claimed}
				if claimed {
					mt.Formed = append(mt.Formed, m.Players[t.Seat].Formed...)
				}
				m.Push(&mt)
			case CMD_RESET:
				m.Reset()
				mt := MahjongResetEvent{Started: false}
				m.Push(&mt)
			}
		}
	}
	tick.Stop()
	close(m.Turn)
	core.AppLog.Printf("table closed %d\n", m.Id)
}

func (m *MahjongTable) Sit(systemId int64, seatNumber int) error {
	if seatNumber < SEAT_E || seatNumber > SEAT_N {
		return fmt.Errorf("invalid seat number %d", seatNumber)
	}
	if m.Players[SEAT_E].SystemId == systemId || m.Players[SEAT_S].SystemId == systemId || m.Players[SEAT_W].SystemId == systemId || m.Players[SEAT_N].SystemId == systemId {
		return fmt.Errorf("already has a seat %d", seatNumber)
	}
	if m.Players[seatNumber].SystemId != 0 {
		return fmt.Errorf("seat already occupied %d", seatNumber)
	}
	m.Players[seatNumber].SystemId = systemId
	m.Players[seatNumber].Auto = false
	return nil

}

func (m *MahjongTable) Dice() []int {
	dice := m.Setup.Dice()
	m.Pts = dice[0] + dice[1]
	return dice
}
func (m *MahjongTable) Deal() error {
	dealer := (m.Pts - 1) % 4
	r := 3
	for {
		if r == 0 {
			m.deal(dealer)
			m.deal(dealer)

		} else {
			m.deal(dealer)
			m.deal(dealer)
			m.deal(dealer)
			m.deal(dealer)
		}
		x := 2
		p := dealer + 1
		for {
			if p == 4 {
				p = 0
			}
			if r == 0 {
				m.deal(p)
			} else {
				m.deal(p)
				m.deal(p)
				m.deal(p)
				m.deal(p)
			}
			p++
			if x == 0 {
				break
			}
			x--
		}
		if r == 0 {
			break
		}
		r--
	}
	m.Started = true
	return nil
}

func (m *MahjongTable) Draw(seat int) error {
	return m.deal(seat)
}
func (m *MahjongTable) Knog(seat int, knog int) error {
	mp := m.Players[seat]
	sz := len(mp.PendingKongs)
	if sz == 0 {
		return fmt.Errorf("no pending kong %d", sz)
	}
	deleted := false
	for i := range mp.PendingKongs {
		if knog == mp.PendingKongs[i] {
			if i == sz-1 {
				mp.PendingKongs = mp.PendingKongs[:sz-1]
			} else {
				mp.PendingKongs = slices.Delete(mp.PendingKongs, i, i+1)
			}
			deleted = true
			break
		}
	}
	if !deleted {
		return fmt.Errorf("no pending kong %d", knog)
	}
	if knog > mj.FS_LIMIT {
		m.Discharge(seat, knog)
	} else {
		//4 kind knog / or pung to knog
	}
	return mp.Knog(&m.Setup.Deck)
}

func (m *MahjongTable) Discharge(seat int, t int) error {
	mp := m.Players[seat]
	sz := len(mp.Tiles)
	if sz == 1 {
		return fmt.Errorf("no more discharge %d", sz)
	}
	drop, err := mp.Hand.Discharge(t)
	if err != nil {
		return err
	}
	m.Discharged = append(m.Discharged, drop)
	return nil
}

func (m *MahjongTable) Chow(seat int, t mj.Tile) error {
	return nil
}

func (m *MahjongTable) Claim(seat int) bool {
	return m.Setup.Mahjong(&m.Players[seat].Hand)
}

func (m *MahjongTable) deal(p int) error {
	mp := m.Players[p]
	sz := len(mp.Tiles)
	if sz > HAND_SIZE_THRESHHOLD {
		return fmt.Errorf("no more draw %d", sz)
	}
	fz := len(mp.Flowers)
	err := mp.Draw(&m.Setup.Deck)
	if err != nil {
		return err
	}
	sz = len(mp.Flowers)
	if fz == sz {
		return nil
	}
	fz = sz
	for {
		err = mp.Knog(&m.Setup.Deck)
		if err != nil {
			return err
		}
		sz = len(mp.Flowers)
		if fz == sz {
			break
		}
		fz = sz
	}
	return nil
}
