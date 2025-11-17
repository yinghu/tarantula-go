package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type MahjongResetEvent struct {
	Started bool
	event.EventObj
}

func (s *MahjongResetEvent) ClassId() int {
	return M_RESET_CID
}

func (s *MahjongResetEvent) ETag() string {
	return "reset"
}

func (s *MahjongResetEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteBool(s.Started); err != nil {
		return err
	}
	return nil
}

func (s *MahjongResetEvent) Outbound(buff core.DataBuffer) error {
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
