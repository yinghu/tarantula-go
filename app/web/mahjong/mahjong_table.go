package main

import (
	"fmt"

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
	CMJ           mj.ClassicMahjong `json:"-"`
	Players       [4]*MahjongPlayer `json:"Players"`
	dice          []int             `json:"-"`
	Discarded     []mj.Tile         `json:"Discharged"`
	Solo          bool
	Turn          chan MahjongPlayToken `json:"-"`
	event.Pusher  `json:"-"`
	core.Sequence `json:"-"`
	Sync          chan MahjongPlayToken `json:"-"`
	skips         []MahjongDiscardEvent

	timerIndex      map[int64]MahjongTimeout
	tp              int
	tpBeforeDiscard int
}

func (m *MahjongTable) Pts() int {
	return m.dice[0] + m.dice[1]
}

func (m *MahjongTable) New() {
	m.CMJ.New()
	m.Players[SEAT_E] = NewPlayer(SEAT_E, true, m)
	m.Players[SEAT_S] = NewPlayer(SEAT_S, true, m)
	m.Players[SEAT_W] = NewPlayer(SEAT_W, true, m)
	m.Players[SEAT_N] = NewPlayer(SEAT_N, true, m)
	m.Discarded = make([]mj.Tile, 0)
	m.Turn = make(chan MahjongPlayToken, 3)
	m.Sync = make(chan MahjongPlayToken, 3)
	m.skips = make([]MahjongDiscardEvent, 0)
	m.timerIndex = map[int64]MahjongTimeout{}
}

func (m *MahjongTable) Reset() {
	m.CMJ.New()
	m.Players[SEAT_E].Reset()
	m.Players[SEAT_S].Reset()
	m.Players[SEAT_W].Reset()
	m.Players[SEAT_N].Reset()
	m.Discarded = m.Discarded[:0]
}

func (m *MahjongTable) Play() {

	running := true

	for running {
		select {
		case k := <-m.Sync: //player control
			switch k.Cmd {
			case CMD_END:
				running = false
				for _, t := range m.timerIndex {
					t.Stop(k, true)
				}
			case CMD_LEAVE:
				m.Leave(k.SystemId)
			case CMD_SIT:
				seat, err := m.Sit(k.SystemId)
				if err != nil {
					mr := NewMahjongErrorEvent(k.SystemId, m.Id, 100, err.Error())
					m.Push(&mr)
				} else {
					m.Players[seat].OnSeat(m)
				}
			}
		case t := <-m.Turn: //play control
			if t.Cmd == CMD_TABLE {
				if m.Solo {
					tm := m.Players[t.Seat].PlayDice(m)
					m.timerIndex[tm.OId()] = tm
					tm.Start(m)
				} else {
					//load table data to players
					tm := m.Players[t.Seat].PlayDice(m)
					m.timerIndex[tm.OId()] = tm
					tm.Start(m)
				}
				continue
			}
			timer, exists := m.timerIndex[t.Id]
			if !exists {
				core.AppLog.Printf("Token not existed: %d selected : %d cmd: %d Id :%d\n", t.Seat, t.Selected, t.Cmd, t.Id)
				m.Players[t.Seat].OnError(m, fmt.Errorf("token not existed %d", t.Id))
				continue
			}
			delete(m.timerIndex, t.Id)
			switch t.Cmd {
			case CMD_DICE:
				m.Dice()
				timer.Stop(t, false)
				tm := m.Players[t.Seat].PlayDeal(m)
				m.timerIndex[tm.OId()] = tm
				tm.Start(m)
			case CMD_DEAL:
				dealer, err := m.Deal()
				if err != nil {
					m.Players[t.Seat].OnError(m, err)
					continue
				}
				m.tp = dealer //start dealer
				timer.Stop(t, false)
				tm := m.Players[m.tp%4].Play(m)
				m.timerIndex[tm.OId()] = tm
				tm.Start(m)
			case CMD_DRAW:
				err := m.Draw(t.Seat)
				if err != nil {
					m.Players[t.Seat].OnError(m, err)
					continue
				}
				m.tp--
				timer.Stop(t, false)
				m.Next()
			case CMD_KONG:
				err := m.Kong(t.Seat, t.Selected)
				if err != nil {
					m.Players[t.Seat].OnError(m, err)
					continue
				}
				m.skips = m.skips[:0]
				m.tp--
				timer.Stop(t, false)
				m.Next()
			case CMD_DISCARD:
				err := m.Discard(t.Seat, t.Selected)
				if err != nil {
					m.Players[t.Seat].OnError(m, err)
					continue
				}
				timer.Stop(t, false)
				m.Next()
			case CMD_CHOW:
				err := m.Chow(t.Seat, t.Selected, t.Chow1, t.Chow2)
				if err != nil {
					m.Players[t.Seat].OnError(m, err)
					continue
				}
				m.skips = m.skips[:0]
				m.tp--
				timer.Stop(t, false)
				m.Next()
			case CMD_PUNG:
				err := m.Pung(t.Seat, t.Selected)
				if err != nil {
					m.Players[t.Seat].OnError(m, err)
					continue
				}
				m.skips = m.skips[:0]
				m.tp--
				timer.Stop(t, false)
				m.Next()
			case CMD_SKIP:
				m.tp = m.tpBeforeDiscard
				m.tp--
				core.AppLog.Printf("skip tp %d %d\n", m.tp, (m.tp+1)%4)
				timer.Stop(t, false)
				m.Next()
			case CMD_CLAIM:
				claimed := m.Claim(t.Seat)
				if !claimed {
					m.Players[t.Seat].OnError(m, fmt.Errorf("fake claimed"))
					continue
				}
				timer.Stop(t, false)
				//m.Reset()
				//go m.Players[t.Seat].PlayDeal(m)
			case CMD_RESET:
				m.Reset()
				mt := MahjongResetEvent{Started: false}
				m.Push(&mt)
				m.Players[t.Seat].PlayDeal(m)
			}
		}
	}
	//time.Sleep(5 * time.Second)
	clear(m.timerIndex)
	close(m.Sync)

	close(m.Turn)
	core.AppLog.Printf("table closed %d\n", m.Id)
}

func (m *MahjongTable) Next() {
	sz := len(m.skips)
	if sz > 0 {
		se := m.skips[sz-1]
		m.skips = m.skips[:sz-1]
		m.tpBeforeDiscard = m.tp
		m.tp = se.Seat
		tm := m.Players[se.Seat].PlayDiscard(m, se)
		m.timerIndex[tm.OId()] = tm
		tm.Start(m)
		return
	}
	m.tp++ //next player turn
	tm := m.Players[m.tp%4].Play(m)
	m.timerIndex[tm.OId()] = tm
	tm.Start(m)
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

func (m *MahjongTable) Leave(systemId int64) {
	if m.Players[SEAT_E].SystemId == systemId {
		m.Players[SEAT_E].SystemId = 0
		m.Players[SEAT_E].Auto = true
		return
	}
	if m.Players[SEAT_S].SystemId == systemId {
		m.Players[SEAT_S].SystemId = 0
		m.Players[SEAT_S].Auto = true
		return
	}
	if m.Players[SEAT_W].SystemId == systemId {
		m.Players[SEAT_W].SystemId = 0
		m.Players[SEAT_W].Auto = true
		return
	}
	if m.Players[SEAT_N].SystemId == systemId {
		m.Players[SEAT_N].SystemId = 0
		m.Players[SEAT_N].Auto = true
		return
	}
}

func (m *MahjongTable) Dice() {
	m.dice = m.CMJ.Dice()
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
	return dealer, nil
}

func (m *MahjongTable) Draw(seat int) error {
	mp := m.Players[seat]
	if mp.TN {
		return fmt.Errorf("cannot draw tile from tn : %v", mp.TN)
	}
	err := m.deal(seat)
	if err != nil {
		return err
	}
	mp.TN = true
	return nil
}
func (m *MahjongTable) Kong(seat int, kong int) error {
	mp := m.Players[seat]
	kt, err := mp.validateKong(kong)
	if err != nil {
		return err
	}
	if kt.Tye == K_FLOWER {
		mp.Discard(mj.FromQ(kong))
		return mp.TailDraw(&m.CMJ.Deck)
	}
	if kt.Tye == K_CONCEALED {
		k := mj.FromQ(kong)
		err = mp.Kong(k)
		if err != nil {
			return err
		}
		return mp.TailDraw(&m.CMJ.Deck)
	}
	if kt.Tye == K_PUNG && !mp.TN {
		k := mj.FromQ(kong)
		err = mp.Kong(k)
		if err != nil {
			return err
		}
		mp.TN = true
		return mp.TailDraw(&m.CMJ.Deck)
	}
	if kt.Tye == K_FORMED && mp.TN {
		core.AppLog.Printf("Rub Kong Checkpoint%v\n", kt)
		p := seat + 1
		claimed1 := m.CMJ.CheckMahjong(&m.Players[p%4].Hand, mj.FromQ(kong))
		p++
		claimed2 := m.CMJ.CheckMahjong(&m.Players[p%4].Hand, mj.FromQ(kong))
		p++
		claimed3 := m.CMJ.CheckMahjong(&m.Players[p%4].Hand, mj.FromQ(kong))
		core.AppLog.Printf("Claimed %v %v %v\n", claimed1, claimed2, claimed3)
		if !claimed1 && !claimed2 && !claimed3 {
			k := mj.FromQ(kong)
			err = mp.Kong(k)
			if err != nil {
				return err
			}
			//exposed kong
			return mp.TailDraw(&m.CMJ.Deck)
		}
		return fmt.Errorf("Rub kong !!!! %v %v %v", claimed1, claimed2, claimed3)
	}
	return fmt.Errorf("no knog from %d", kong)
}

func (m *MahjongTable) Discard(seat int, t int) error {
	mp := m.Players[seat]
	if !mp.TN {
		return fmt.Errorf("cannot discard tile from tn : %v", mp.TN)
	}
	drop := mj.FromQ(t)
	err := mp.Discard(drop)
	if err != nil {
		return err
	}
	m.Discarded = append(m.Discarded, drop)
	mp.TN = false
	p := seat + 1
	pp := m.Players[p%4].CheckDiscard(seat, drop, true) //pung/kong/chow
	if len(pp) > 0 {
		m.skips = append(m.skips, NewMahjongDiscardEvent(p%4, seat, drop, pp))
	}
	p++
	pp1 := m.Players[p%4].CheckDiscard(seat, drop, false) //pung/kong
	if len(pp1) > 0 {
		m.skips = append(m.skips, NewMahjongDiscardEvent(p%4, seat, drop, pp1))
		return nil
	}
	p++
	pp2 := m.Players[p%4].CheckDiscard(seat, drop, false) //pung/kong
	if len(pp2) > 0 {
		m.skips = append(m.skips, NewMahjongDiscardEvent(p%4, seat, drop, pp2))
	}
	return nil
}

func (m *MahjongTable) Chow(seat int, drop int, chow1 int, chow2 int) error {
	mp := m.Players[seat]
	if mp.TN {
		return fmt.Errorf("cannot chow tile from tn : %v", mp.TN)
	}
	d := mj.FromQ(drop)
	c1 := mj.FromQ(chow1)
	c2 := mj.FromQ(chow2)
	err := mp.Chow(d, c1, c2)
	if err != nil {
		return err
	}
	mp.TN = true
	return nil
}

func (m *MahjongTable) Pung(seat int, drop int) error {
	mp := m.Players[seat]
	if mp.TN {
		return fmt.Errorf("cannot pung tile from tn : %v", mp.TN)
	}
	d := mj.FromQ(drop)
	err := mp.Pung(d)
	if err != nil {
		return err
	}
	mp.TN = true
	return nil
}

func (m *MahjongTable) Claim(seat int) bool {
	mp := m.Players[seat]
	if !mp.TN {
		return false
	}
	return m.CMJ.Mahjong(&mp.Hand)
}

func (m *MahjongTable) deal(p int) error {
	mp := m.Players[p]
	sz := len(mp.Tiles)
	if sz > HAND_SIZE_THRESHHOLD {
		return fmt.Errorf("no more draw %d", sz)
	}
	fz := len(mp.Flowers)
	err := mp.HeadDraw(&m.CMJ.Deck)
	if err != nil {
		return err
	}
	sz = len(mp.Flowers)
	if fz == sz {
		return nil
	}
	fz = sz
	for {
		err = mp.TailDraw(&m.CMJ.Deck)
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
