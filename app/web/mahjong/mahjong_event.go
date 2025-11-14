package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

const(
	M_TOKEN_CID int = 100
	M_DICE_CID int = 101
	M_HAND_CID int = 102
	M_TABLE_CID int = 103
	M_SIT_CID int = 104
)

type MahjongEvent struct {
	Token    MahjongPlayToken
	SystemId int64
	event.EventObj
}

func (s *MahjongEvent) ClassId() int {
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
	return s.Token.Read(buff)
}

func (s *MahjongEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(s.SystemId); err != nil {
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

func (s *MahjongEvent) RecipientId() int64 {
	return s.SystemId
}
