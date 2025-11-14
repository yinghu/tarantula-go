package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type MahjongHandEvent struct {
	Table *MahjongTable
	event.EventObj
}

func (s *MahjongHandEvent) ClassId() int {
	return 102
}

func (s *MahjongHandEvent) ETag() string {
	return "table"
}

func (s *MahjongHandEvent) Write(buff core.DataBuffer) error {
	return s.Table.Players[SEAT_E].Hand.Write(buff)
}

func (s *MahjongHandEvent) Outbound(buff core.DataBuffer) error {
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
