package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type MahjongErrorEvent struct {
	SystemId int64
	TableId  int64
	Code     int
	Message  string
	event.EventObj
}

func (s *MahjongErrorEvent) ClassId() int {
	return M_ERR_CID
}

func (s *MahjongErrorEvent) ETag() string {
	return "error"
}

func (s *MahjongErrorEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(s.SystemId); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.TableId); err != nil {
		return err
	}
	if err := buff.WriteInt32(int32(s.Code)); err != nil {
		return err
	}
	if err := buff.WriteString(s.Message); err != nil {
		return err
	}
	return nil
}

func (s *MahjongErrorEvent) Outbound(buff core.DataBuffer) error {
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
