package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/mj"
)

type MahjongDischargeEvent struct {
	Seat     int
	DropSeat int
	Drop     mj.Tile
	Opts     []mj.Meld
	MahjongEventObj
}

func (s *MahjongDischargeEvent) ClassId() int {
	return M_DISCHARGE_CID
}

func (s *MahjongDischargeEvent) ETag() string {
	return "discharge"
}

func (s *MahjongDischargeEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt32(int32(s.Seat)); err != nil {
		return err
	}
	if err := buff.WriteInt32(int32(s.DropSeat)); err != nil {
		return err
	}
	if err := s.Drop.Write(buff); err!=nil{
		return err
	}
	for _, m := range s.Opts {
		if err := m.Write(buff); err != nil {
			return err
		}
	}
	return nil
}

func (s *MahjongDischargeEvent) Outbound(buff core.DataBuffer) error {
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
