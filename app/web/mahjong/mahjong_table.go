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

	HAND_SIZE_THRESHHOLD       int = 13
	PLAYER_TURN_TICKER_SECONDS int = 30
	AUTO_TURN_TICKER_SECONDS   int = 3
	COUNT_DOWN_BUFFER          int = 1
	SOLO                       int = 1
)

type MahjongTable struct {
	Id            int64             `json:"Id,string"`
	Setup         mj.ClassicMahjong `json:"-"`
	Players       [4]*MahjongPlayer `json:"Players"`
	dice          []int             `json:"-"`
	Discharged    []mj.Tile         `json:"Discharged"`
	Started       bool
	Solo          bool
	Turn          chan MahjongPlayToken `json:"-"`
	event.Pusher  `json:"-"`
	core.Sequence `json:"-"`
	Timer         chan MahjongTimeout
	Sync          chan MahjongPlayToken `json:"-"`
	skips         []MahjongDischargeEvent
}

func (m *MahjongTable) Pts() int {
	return m.dice[0] + m.dice[1]
}

func (m *MahjongTable) New() {
	m.Setup.New()
	m.Players[SEAT_E] = NewPlayer(SEAT_E, true, m)
	m.Players[SEAT_S] = NewPlayer(SEAT_S, true, m)
	m.Players[SEAT_W] = NewPlayer(SEAT_W, true, m)
	m.Players[SEAT_N] = NewPlayer(SEAT_N, true, m)
	m.Discharged = make([]mj.Tile, 0)
	m.Turn = make(chan MahjongPlayToken, 3)
	m.Timer = make(chan MahjongTimeout, 3)
	m.Sync = make(chan MahjongPlayToken, 3)
	m.skips = make([]MahjongDischargeEvent, 3)
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
	timerIndex := make(map[int64]MahjongTimeout)
	running := true
	var tp int
	for running {
		select {
		case k := <-m.Sync:
			//core.AppLog.Printf("Sync: %d selected : %d cmd: %d Id :%d\n", k.Seat, k.Selected, k.Cmd, k.Id)
			switch k.Cmd {
			case CMD_END:
				m.Started = false
				running = false
			case CMD_SIT:
				seat, err := m.Sit(k.SystemId)
				if err != nil {
					mr := NewMahjongErrorEvent(k.SystemId, m.Id, 100, err.Error())
					m.Push(&mr)
				} else {
					go m.Players[seat].OnSeat(m)
				}
			case CMD_TURN_CONTINUE:
				go m.Players[tp%4].Play(m)
			case CMD_TURN_END:
				sz := len(m.skips)
				if sz > 0 {
					se := m.skips[sz-1]
					m.skips = m.skips[:sz-1]
					core.AppLog.Printf("discharge %d %v\n",sz,se.Opts)
					go m.Players[se.Seat].PlayDischarge(m, se)
					break
				}
				tp++ //next player turn
				go m.Players[tp%4].Play(m)
			}
		case tm := <-m.Timer:
			timerIndex[tm.OId()] = tm
			go tm.Start(m)
		case t := <-m.Turn:
			if t.Cmd == CMD_TABLE {
				if m.Solo {
					go m.Players[t.Seat].PlayDice(m)
				} else {
					//load table data to players
				}
				continue
			}
			timer, exists := timerIndex[t.Id]
			if !exists {
				core.AppLog.Printf("Token not existed: %d selected : %d cmd: %d Id :%d\n", t.Seat, t.Selected, t.Cmd, t.Id)
				continue
			}
			delete(timerIndex, t.Id)
			switch t.Cmd {
			case CMD_DICE:
				m.Dice()
				timer.Stop(true, false)
				go m.Players[t.Seat].PlayDeal(m)
			case CMD_DEAL:
				dealer, err := m.Deal()
				if err != nil {
					go m.Players[t.Seat].OnError(m, err)
					continue
				}
				tp = dealer
				timer.Stop(true, false)
				//start dealer
			case CMD_DRAW:
				err := m.Draw(t.Seat)
				if err != nil {
					go m.Players[t.Seat].OnError(m, err)
					continue
				}
				timer.Stop(true, false)
			case CMD_KONG:
				err := m.Kong(t.Seat, t.Selected)
				if err != nil {
					go m.Players[t.Seat].OnError(m, err)
					continue
				}
				timer.Stop(true, false)
			case CMD_DISCHARGE:
				err := m.Discharge(t.Seat, t.Selected)
				if err != nil {
					go m.Players[t.Seat].OnError(m, err)
					continue
				}
				timer.Stop(true, false)
			case CMD_CHOW:
				err := m.Chow(t.Seat, t.Selected, t.Chow1, t.Chow2)
				if err != nil {
					go m.Players[t.Seat].OnError(m, err)
					continue
				}
				tp = t.Seat
				timer.Stop(true, false)
			case CMD_PUNG:
				err := m.Pung(t.Seat, t.Selected)
				if err != nil {
					go m.Players[t.Seat].OnError(m, err)
					continue
				}
				tp = t.Seat
				timer.Stop(true, false)
			case CMD_SKIP:
				timer.Stop(false, false)
			case CMD_CLAIM:
				claimed := m.Claim(t.Seat)
				if !claimed {
					go m.Players[t.Seat].OnError(m, fmt.Errorf("fake claimed"))
					continue
				}
				timer.Stop(true, false)
				m.Reset()
				go m.Players[t.Seat].PlayDeal(m)
			case CMD_RESET:
				m.Reset()
				mt := MahjongResetEvent{Started: false}
				m.Push(&mt)
				go m.Players[t.Seat].PlayDeal(m)
			}
		}
	}
	for _, t := range timerIndex {
		t.Stop(false, true)
	}
	time.Sleep(5 * time.Second)
	clear(timerIndex)
	close(m.Sync)
	close(m.Timer)
	close(m.Turn)
	core.AppLog.Printf("table closed %d\n", m.Id)
}

func (m *MahjongTable) Sit(systemId int64) (int, error) {
	if m.Players[SEAT_E].SystemId == 0 {
		m.Players[SEAT_E].SystemId = systemId
		m.Players[SEAT_E].Auto = false
		return SEAT_E, nil
	}
	if m.Players[SEAT_S].SystemId == 0 {
		m.Players[SEAT_S].SystemId = systemId
		m.Players[SEAT_S].Auto = false
		return SEAT_S, nil
	}
	if m.Players[SEAT_W].SystemId == 0 {
		m.Players[SEAT_W].SystemId = systemId
		m.Players[SEAT_W].Auto = false
		return SEAT_W, nil
	}
	if m.Players[SEAT_N].SystemId == 0 {
		m.Players[SEAT_N].SystemId = systemId
		m.Players[SEAT_N].Auto = false
		return SEAT_N, nil
	}
	return 0, fmt.Errorf("no seat available on table: %d", m.Id)

}

func (m *MahjongTable) Dice() {
	m.dice = m.Setup.Dice()
}
func (m *MahjongTable) Deal() (int, error) {
	if m.Pts() == 0 {
		return 0, fmt.Errorf("no dice")
	}
	dealer := (m.Pts() - 1) % 4
	m.Players[dealer].TN = true
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
	return dealer, nil
}

func (m *MahjongTable) Draw(seat int) error {
	return m.deal(seat)
}
func (m *MahjongTable) Kong(seat int, kong int) error {
	mp := m.Players[seat]
	sz := len(mp.PendingKongs)
	if sz == 0 {
		return fmt.Errorf("no pending kong %d", sz)
	}
	deleted := false
	for i := range mp.PendingKongs {
		if kong == mp.PendingKongs[i] {
			mp.PendingKongs = slices.Delete(mp.PendingKongs, i, i+1)
			deleted = true
			break
		}
	}
	if !deleted {
		return fmt.Errorf("no pending kong %d", kong)
	}
	if kong > mj.FS_LIMIT {
		m.Discharge(seat, kong)
		return mp.TailDraw(&m.Setup.Deck)
	}
	k := mj.FromQ(kong)
	err := mp.Kong(k)
	if err != nil {
		return err
	}
	//4 kind kong / or pung to kong
	return mp.TailDraw(&m.Setup.Deck)
}

func (m *MahjongTable) Discharge(seat int, t int) error {
	mp := m.Players[seat]
	sz := len(mp.Tiles)
	if sz == 1 {
		return fmt.Errorf("no more discharge %d", sz)
	}
	drop := mj.FromQ(t)
	err := mp.Hand.Discharge(drop)
	if err != nil {
		return err
	}
	m.Discharged = append(m.Discharged, drop)
	p := seat + 1
	pp := m.Players[p%4].CheckDischarge(seat, drop, true) //pung/kong/chow
	if len(pp) > 0 {
		core.AppLog.Printf("pp : %v\n", pp)
		m.skips = append(m.skips, NewMahjongDischargeEvent(p%4, seat, drop, pp))
	}
	p++
	pp1 := m.Players[p%4].CheckDischarge(seat, drop, false) //pung/kong
	if len(pp1) > 0 {
		core.AppLog.Printf("1PP : %v\n", pp1)
		m.skips = append(m.skips, NewMahjongDischargeEvent(p%4, seat, drop, pp1))
		return nil
	}
	p++
	pp2 := m.Players[p%4].CheckDischarge(seat, drop, false) //pung/kong
	if len(pp2) > 0 {
		core.AppLog.Printf("2PP : %v\n", pp2)
		m.skips = append(m.skips, NewMahjongDischargeEvent(p%4, seat, drop, pp2))
	}
	return nil
}

func (m *MahjongTable) Chow(seat int, drop int, chow1 int, chow2 int) error {
	d := mj.FromQ(drop)
	c1 := mj.FromQ(chow1)
	c2 := mj.FromQ(chow2)
	mp := m.Players[seat]
	return mp.Chow(d, c1, c2)
}

func (m *MahjongTable) Pung(seat int, drop int) error {
	d := mj.FromQ(drop)
	mp := m.Players[seat]
	return mp.Pung(d)
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
	err := mp.HeadDraw(&m.Setup.Deck)
	if err != nil {
		return err
	}
	sz = len(mp.Flowers)
	if fz == sz {
		return nil
	}
	fz = sz
	for {
		err = mp.TailDraw(&m.Setup.Deck)
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

func (m *MahjongTable) Update(e event.Event) {
	if e.RecipientId() > 0 {
		m.Push(e)
		return
	}
	//table broadcasting
	for i := range m.Players {
		if m.Players[i].SystemId > 0 {
			e.OnRecipientId(m.Players[i].SystemId)
			m.Push(e)
		}
	}
}
