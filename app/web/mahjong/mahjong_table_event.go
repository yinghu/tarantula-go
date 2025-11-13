package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type MahjongTableEvent struct {
	Table *MahjongTable
	event.EventObj
}

func (s *MahjongTableEvent) ClassId() int {
	return 102
}

func (s *MahjongTableEvent) ETag() string {
	return "table"
}

func (s *MahjongTableEvent) Write(buff core.DataBuffer) error {
	return s.Table.Players[SEAT_E].Hand.Write(buff)
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
