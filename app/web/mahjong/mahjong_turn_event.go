package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type MahjongTurnEvent struct {
	SystemId int64
	TableId  int64
	Seat     int32
	event.EventObj
}

func (s *MahjongTurnEvent) ClassId() int {
	return M_TURN_CID
}

func (s *MahjongTurnEvent) ETag() string {
	return "turn"
}

func (s *MahjongTurnEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(s.SystemId); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.TableId); err != nil {
		return err
	}
	if err := buff.WriteInt32(s.Seat); err != nil {
		return err
	}
	return nil
}

func (s *MahjongTurnEvent) Outbound(buff core.DataBuffer) error {
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
