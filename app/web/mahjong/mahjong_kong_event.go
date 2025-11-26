package main

import (
	"gameclustering.com/internal/core"
)

type MahjongKongEvent struct {
	
	Knog     []int
	MahjongEventObj
}

func (s *MahjongKongEvent) ClassId() int {
	return M_KNOG_CID
}

func (s *MahjongKongEvent) ETag() string {
	return "kong"
}

func (s *MahjongKongEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(s.SystemId); err != nil {
		return err
	}
	sz := len(s.Knog)
	if err := buff.WriteInt32(int32(sz)); err != nil {
		return err
	}
	for i := range s.Knog {
		if err := buff.WriteInt32(int32(s.Knog[i])); err != nil {
			return err
		}
	}
	return nil
}

func (s *MahjongKongEvent) Outbound(buff core.DataBuffer) error {
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

