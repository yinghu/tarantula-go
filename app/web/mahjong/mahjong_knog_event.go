package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type MahjongKnogEvent struct {
	SystemId int64
	Knog     int
	event.EventObj
}

func (s *MahjongKnogEvent) ClassId() int {
	return M_KNOG_CID
}

func (s *MahjongKnogEvent) ETag() string {
	return "knog"
}

func (s *MahjongKnogEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(s.SystemId); err != nil {
		return err
	}
	if err := buff.WriteInt32(int32(s.Knog)); err != nil {
		return err
	}
	return nil
}

func (s *MahjongKnogEvent) Outbound(buff core.DataBuffer) error {
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

func (s *MahjongKnogEvent) RecipientId() int64 {
	return s.SystemId
}
