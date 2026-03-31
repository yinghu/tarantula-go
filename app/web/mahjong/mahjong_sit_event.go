package main

import (
	"gameclustering.com/internal/core"
)

type MahjongSitEvent struct {
	Seat     int32
	MahjongEventObj
}

func (s *MahjongSitEvent) ClassId() uint32 {
	return M_SIT_CID
}

func (s *MahjongSitEvent) ETag() string {
	return "sit"
}

func (s *MahjongSitEvent) Write(buff core.DataBuffer) error {
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

func (s *MahjongSitEvent) Outbound(buff core.DataBuffer) error {
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
