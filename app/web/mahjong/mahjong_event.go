package main

import (
	"gameclustering.com/internal/core"
)

const (
	M_TOKEN_CID     int32 = 100
	M_DICE_CID      int32 = 101
	M_HAND_CID      int32 = 102
	M_TABLE_CID     int32 = 103
	M_SIT_CID       int32 = 104
	M_CLAIM_CID     int32 = 105
	M_RESET_CID     int32 = 106
	M_DISCHARGE_CID int32 = 107
	M_KNOG_CID      int32 = 108
	M_TURN_CID      int32 = 109
	M_ERR_CID       int32 = 911
)

type MahjongEvent struct {
	Token MahjongPlayToken
	MahjongEventObj
}

func (s *MahjongEvent) ClassId() int32 {
	return M_TOKEN_CID
}

func (s *MahjongEvent) ETag() string {
	return "mj"
}

func (s *MahjongEvent) WriteKey(buff core.DataBuffer) error {
	if err := buff.WriteString(s.ETag()); err != nil {
		return err
	}
	return nil
}

func (s *MahjongEvent) ReadKey(buff core.DataBuffer) error {
	_, err := buff.ReadString()
	if err != nil {
		return err
	}
	return nil
}

func (s *MahjongEvent) Read(buff core.DataBuffer) error {
	sysId, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.SystemId = sysId
	tableId, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.TableId = tableId

	return s.Token.Read(buff)
}

func (s *MahjongEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(s.SystemId); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.TableId); err != nil {
		return err
	}
	return s.Token.Write(buff)
}

func (s *MahjongEvent) Outbound(buff core.DataBuffer) error {
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

func (s *MahjongEvent) Inbound(buff core.DataBuffer) error {
	err := s.ReadKey(buff)
	if err != nil {
		s.Callback.OnError(s, err)
		return err
	}
	err = s.Read(buff)
	if err != nil {
		s.Callback.OnError(s, err)
		return err
	}
	s.Callback.OnEvent(s)
	return nil
}
