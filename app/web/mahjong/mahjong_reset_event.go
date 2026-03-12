package main

import (
	"gameclustering.com/internal/core"
)

type MahjongResetEvent struct {
	Started bool
	MahjongEventObj
}

func (s *MahjongResetEvent) ClassId() uint32 {
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
