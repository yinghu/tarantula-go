package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/mj"
)

type MahjongHandEvent struct {
	H mj.Hand
	K []int
	event.EventObj
}

func (s *MahjongHandEvent) ClassId() int {
	return M_HAND_CID
}

func (s *MahjongHandEvent) ETag() string {
	return "hand"
}

func (s *MahjongHandEvent) Write(buff core.DataBuffer) error {
	err := s.H.Write(buff)
	if err != nil {
		return err
	}
	sz := len(s.K)
	if err := buff.WriteInt32(int32(sz)); err != nil {
		return err
	}
	for i := range s.K {
		if err := buff.WriteInt32(int32(s.K[i])); err != nil {
			return err
		}
	}
	return nil
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
