package main

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/mj"
)

const (
	SEAT_E int = 0
	SEAT_S int = 1
	SEAT_W int = 2
	SEAT_N int = 3

	HAND_SIZE_THRESHHOLD int = 13
)

type MahjongTable struct {
	Id              int64             `json:"Id,string"`
	Setup           mj.ClassicMahjong `json:"-"`
	Players         [4]*MahjongPlayer `json:"Players"`
	Pts             int               `json:"Pts"`
	Discharged      []mj.Tile         `json:"Discharged"`
	Started         bool
	Turn            chan MahjongPlayToken `json:"-"`
	*MahjongService `json:"-"`
}

func (m *MahjongTable) Reset() {
	m.Setup.New()
	m.Players[SEAT_E] = NewPlayer("East")
	m.Players[SEAT_S] = NewPlayer("South")
	m.Players[SEAT_W] = NewPlayer("West")
	m.Players[SEAT_N] = NewPlayer("North")
	m.Discharged = make([]mj.Tile, 0)
}

func (m *MahjongTable) Play() {
	for t := range m.Turn {
		m.Reset()
		m.Dice()
		m.Deal()
		switch t.Cmd {
		case CMD_SIT:
			err := m.Sit(t.SystemId, t.Seat)
			if err != nil {
				mr := MahjongErrorEvent{SystemId: t.SystemId, TableId: m.Id, Code: 100, Message: err.Error()}
				m.MahjongService.Pusher().Push(&mr)
			} else {
				mt := MahjongSitEvent{SystemId: t.SystemId, TableId: m.Id}
				m.MahjongService.Pusher().Push(&mt)
			}
		case CMD_DICE:
			mt := MahjongHandEvent{Table: m}
			m.MahjongService.Pusher().Push(&mt)
		case CMD_DEAL:
			mt := MahjongHandEvent{Table: m}
			m.MahjongService.Pusher().Push(&mt)
		case CMD_END:
			return
		}
	}
}

func (m *MahjongTable) Sit(systemId int64, seatNumber int) error {
	core.AppLog.Printf("Sit : %d > %d >%d", systemId, seatNumber, m.Players[seatNumber].SystemId)
	switch seatNumber {
	case SEAT_E:
		if m.Players[SEAT_E].SystemId != 0 {
			return fmt.Errorf("seat already occupied %d", seatNumber)
		}
		m.Players[SEAT_E].SystemId = systemId
		m.Players[SEAT_E].Auto = false
		return nil
	case SEAT_S:
		if m.Players[SEAT_S].SystemId != 0 {
			return fmt.Errorf("seat already occupied %d", seatNumber)
		}
		m.Players[SEAT_S].SystemId = systemId
		m.Players[SEAT_S].Auto = false
		return nil
	case SEAT_W:
		if m.Players[SEAT_W].SystemId != 0 {
			return fmt.Errorf("seat already occupied %d", seatNumber)
		}
		m.Players[SEAT_W].SystemId = systemId
		m.Players[SEAT_W].Auto = false
		return nil
	case SEAT_N:
		if m.Players[SEAT_N].SystemId != 0 {
			return fmt.Errorf("seat already occupied %d", seatNumber)
		}
		m.Players[SEAT_N].SystemId = systemId
		m.Players[SEAT_N].Auto = false
		return nil
	}
	return fmt.Errorf("invalid seat number %d", seatNumber)
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
	mp := m.Players[seat]
	sz := len(mp.Tiles)
	if sz > HAND_SIZE_THRESHHOLD {
		return fmt.Errorf("no more draw %d", sz)
	}
	return m.deal(seat)
}

func (m *MahjongTable) Discharge(seat int, t mj.Tile) error {
	mp := m.Players[seat]
	sz := len(mp.Tiles)
	if sz == 1 {
		return fmt.Errorf("no more discharge %d", sz)
	}
	err := mp.Drop(t)
	if err != nil {
		return err
	}
	m.Players[seat] = mp
	m.Discharged = append(m.Discharged, t)
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
	fz := len(mp.Flowers)
	err := mp.Draw(&m.Setup.Deck)
	if err != nil {
		return err
	}
	sz := len(mp.Flowers)
	if fz == sz {
		m.Players[p] = mp
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
	m.Players[p] = mp
	return nil
}
