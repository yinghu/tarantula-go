package main

import (
	"time"

	"gameclustering.com/internal/core"
)

type MahjongTableEvent struct {
	SystemId int64
	TableId  int64
	East     int64
	South    int64
	West     int64
	North    int64
	MahjongTimeoutObj
}

func (s *MahjongTableEvent) ClassId() int {
	return M_TABLE_CID
}

func (s *MahjongTableEvent) ETag() string {
	return "table"
}

func (s *MahjongTableEvent) WriteKey(buff core.DataBuffer) error {
	if err := buff.WriteInt64(s.OId()); err != nil {
		return err
	}
	return nil
}

func (s *MahjongTableEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(s.SystemId); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.TableId); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.East); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.South); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.West); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.North); err != nil {
		return err
	}
	return s.Next.Write(buff)
}

func (s *MahjongTableEvent) Outbound(buff core.DataBuffer) error {
	err := s.WriteKey(buff)
	if err != nil {
		s.Callback.OnError(s, err)
		return err
	}
	err = s.Write(buff)
	if err != nil {
		s.Callback.OnError(s, err)
		return err
	}
	return nil
}

func (s *MahjongTableEvent) Start(tb *MahjongTable) {
	tm := *time.NewTimer(time.Duration(s.Next.CountDown+5) * time.Second)
	s.Commited = make(chan bool)
	closing := false
	defer close(s.Commited)
	for {
		if closing {
			break
		}
		select {
		case <-tm.C:
			se := MahjongPlayToken{SystemId: s.SystemId, Cmd: CMD_SIT, Seat: SEAT_E, Id: s.OId()}
			tb.Turn <- se
		case <-s.Commited:
			core.AppLog.Printf("Stop %d\n", s.OId())
			closing = true
		}
	}
}
func (s *MahjongTableEvent) Stop() {
	s.Commited <- true
}
