package main

import (
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type MahjongTableEvent struct {
	SystemId  int64
	TableId   int64
	East      int64
	South     int64
	West      int64
	North     int64
	CountDown int
	event.EventObj
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
	if err := buff.WriteInt32(int32(s.CountDown)); err != nil {
		return err
	}
	return nil
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

func (s *MahjongTableEvent) Timeout() {
	tm := time.NewTimer(10 * time.Second)
	t := <-tm.C
	tm.Stop()
	core.AppLog.Printf("Timeout %v\n", t)
}
