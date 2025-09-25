package main

import (
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/core"
)

type MahjongEvent struct {
	Cmd string
	event.EventObj
}

func (s *MahjongEvent) ClassId() int {
	return 1
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
	cmd, err := buff.ReadString()
	if err!=nil {
		return err
	}
	s.Cmd = cmd;
	return nil
}

func (s *MahjongEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteString(s.ETag()); err != nil {
		return err
	}
	return nil
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
