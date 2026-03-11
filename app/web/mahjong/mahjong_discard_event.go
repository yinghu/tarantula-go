package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/mj"
)

type MahjongDiscardEvent struct {
	Seat     int
	DropSeat int
	Drop     mj.Tile
	Opts     []mj.Meld
	MahjongEventObj
}

func (s *MahjongDiscardEvent) ClassId() int32 {
	return M_DISCHARGE_CID
}

func (s *MahjongDiscardEvent) ETag() string {
	return "discharge"
}

func (s *MahjongDiscardEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt32(int32(s.Seat)); err != nil {
		return err
	}
	if err := buff.WriteInt32(int32(s.DropSeat)); err != nil {
		return err
	}
	if err := s.Drop.Write(buff); err != nil {
		return err
	}
	sz := len(s.Opts)
	if err := buff.WriteInt32(int32(sz)); err != nil {
		return err
	}
	for _, m := range s.Opts {
		if err := m.Write(buff); err != nil {
			return err
		}
	}
	return nil
}

func (s *MahjongDiscardEvent) Outbound(buff core.DataBuffer) error {
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
