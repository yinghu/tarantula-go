package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type MahjongClaimEvent struct {
	SystemId int64
	TableId  int64
	Seat     int32
	Claimed  bool
	event.EventObj
}

func (s *MahjongClaimEvent) ClassId() int {
	return M_CLAIM_CID
}

func (s *MahjongClaimEvent) ETag() string {
	return "claim"
}

func (s *MahjongClaimEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(s.SystemId); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.TableId); err != nil {
		return err
	}
	if err := buff.WriteInt32(s.Seat); err != nil {
		return err
	}
	if err := buff.WriteBool(s.Claimed); err != nil {
		return err
	}
	return nil
}

func (s *MahjongClaimEvent) Outbound(buff core.DataBuffer) error {
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
