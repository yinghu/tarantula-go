package main

import (
	"time"

	"gameclustering.com/internal/core"
)

type MahjongTableEvent struct {
	SystemId  int64
	TableId   int64
	East      int64
	South     int64
	West      int64
	North     int64
	CountDown int
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

func (s *MahjongTableEvent) Start(tb *MahjongTable) {
	tm := *time.NewTimer(time.Duration(s.CountDown+5) * time.Second)
	t := <-tm.C
	if s.Commited {
		return
	}
	core.AppLog.Printf("Timeout %v %v\n", t, s.Commited)
	se := MahjongPlayToken{SystemId: s.SystemId, Cmd: CMD_SIT, Seat: SEAT_E, Id: s.OId()}
	tb.Turn <- se
	//mr := MahjongErrorEvent{SystemId: s.SystemId, TableId: s.TableId, Code: 100, Message: "timeout"}
	//tb.Push(&mr)
	//ko := event.KickoffEvent{SystemId: s.SystemId, Source: "timeout"}
	//tb.Push(&ko)

}
func (s *MahjongTableEvent) Stop() {
	s.Commited = true
	core.AppLog.Printf("Stop %v\n", s.Commited)
}
