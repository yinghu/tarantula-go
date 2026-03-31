package main

import (
	"gameclustering.com/internal/core"
)

type MahjongDiceEvent struct {
	Dice1 int32
	Dice2 int32
	MahjongEventObj
}

func (s *MahjongDiceEvent) ClassId() uint32 {
	return M_DICE_CID
}

func (s *MahjongDiceEvent) ETag() string {
	return "dice"
}

func (s *MahjongDiceEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt32(s.Dice1); err != nil {
		return err
	}
	if err := buff.WriteInt32(s.Dice2); err != nil {
		return err
	}
	return nil
}

func (s *MahjongDiceEvent) Outbound(buff core.DataBuffer) error {
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
