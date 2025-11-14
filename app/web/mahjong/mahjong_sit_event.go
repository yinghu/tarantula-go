package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type MahjongSitEvent struct {
	TableId int64
	event.EventObj
}
func (s *MahjongSitEvent) ClassId() int {
	return 101
}

func (s *MahjongSitEvent) ETag() string {
	return "sit"
}

func (s *MahjongSitEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(s.TableId); err != nil {
		return err
	}
	return nil
}


func (s *MahjongSitEvent) Outbound(buff core.DataBuffer) error {
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
