package main

import (
	"gameclustering.com/internal/core"
)

type MahjongTableEvent struct {
	East     int64
	South    int64
	West     int64
	North    int64
	MahjongTimeoutObj
}

func (s *MahjongTableEvent) ClassId() int32 {
	return M_TABLE_CID
}

func (s *MahjongTableEvent) ETag() string {
	return "table"
}

func (s *MahjongTableEvent) WriteKey(buff core.DataBuffer) error {
	if err := buff.WriteInt64(s.OId()); err != nil {
		return err
	}
	return nil
}

func (s *MahjongTableEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(s.SystemId); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.TableId); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.East); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.South); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.West); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.North); err != nil {
		return err
	}
	return s.N.Write(buff)
}

func (s *MahjongTableEvent) Outbound(buff core.DataBuffer) error {
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

