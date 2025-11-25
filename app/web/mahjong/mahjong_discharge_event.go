package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/mj"
)

type MahjongDischargeEvent struct {
	D []mj.Tile
	MahjongEventObj
}

func (s *MahjongDischargeEvent) ClassId() int {
	return M_DISCHARGE_CID
}

func (s *MahjongDischargeEvent) ETag() string {
	return "discharge"
}

func (s *MahjongDischargeEvent) Write(buff core.DataBuffer) error {

	sz := len(s.D)
	if err := buff.WriteInt32(int32(sz)); err != nil {
		return err
	}
	for i := range s.D {
		if err := s.D[i].Write(buff); err != nil {
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
