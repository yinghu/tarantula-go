package main

import (
	"fmt"

	"gameclustering.com/internal/mj"
)

const (
	SEAT_E = 0
	SEAT_S = 1
	SEAT_W = 2
	SEAT_N = 3
)

type MahjongTable struct {
	Setup   mj.ClassicMahjong
	Players [4]MahjongPlayer
	Pts     int
}

func (m *MahjongTable) New() {
	m.Setup.New()
	m.Players[SEAT_E] = NewPlayer()
	m.Players[SEAT_S] = NewPlayer()
	m.Players[SEAT_W] = NewPlayer()
	m.Players[SEAT_N] = NewPlayer()
}

func (m *MahjongTable) Sit(systemId int64, seatNumber int) error {
	switch seatNumber {
	case SEAT_E:
		if m.Players[SEAT_E].SystemId != 0 {
			return fmt.Errorf("seat already occupied %d", seatNumber)
		}
		m.Players[SEAT_E].SystemId = systemId
		return nil
	case SEAT_S:
		if m.Players[SEAT_S].SystemId != 0 {
			return fmt.Errorf("seat already occupied %d", seatNumber)
		}
		m.Players[SEAT_S].SystemId = systemId
		return nil
	case SEAT_W:
		if m.Players[SEAT_W].SystemId != 0 {
			return fmt.Errorf("seat already occupied %d", seatNumber)
		}
		m.Players[SEAT_W].SystemId = systemId
		return nil
	case SEAT_N:
		if m.Players[SEAT_N].SystemId != 0 {
			return fmt.Errorf("seat already occupied %d", seatNumber)
		}
		m.Players[SEAT_N].SystemId = systemId
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
			m.Players[dealer].Draw(&m.Setup.Deck)
			m.Players[dealer].Draw(&m.Setup.Deck)

		} else {
			m.Players[dealer].Draw(&m.Setup.Deck)
			m.Players[dealer].Draw(&m.Setup.Deck)
			m.Players[dealer].Draw(&m.Setup.Deck)
			m.Players[dealer].Draw(&m.Setup.Deck)
		}
		x := 2
		p := dealer + 1
		for {
			if p == 4 {
				p = 0
			}
			if r == 0 {
				err := m.Players[p].Draw(&m.Setup.Deck)
				if err != nil {
					return err
				}
			} else {
				err := m.Players[p].Draw(&m.Setup.Deck)
				if err != nil {
					return err
				}
				err = m.Players[p].Draw(&m.Setup.Deck)
				if err != nil {
					return err
				}
				err = m.Players[p].Draw(&m.Setup.Deck)
				if err != nil {
					return err
				}
				err = m.Players[p].Draw(&m.Setup.Deck)
				if err != nil {
					return err
				}
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
	return nil
}

func (m *MahjongTable) Discharge() error {
	return nil
}

func (m *MahjongTable) Chow() error {
	return nil
}
