package event

import (
	"errors"

	"gameclustering.com/internal/core"
)

type StatEvent struct {
	Name     string `json:"Name"`
	Count    uint32 `json:"Count"`
	core.EventObj `json:"-"`
}

func (s *StatEvent) ClassId() uint32 {
	return STAT_CID
}



func (s *StatEvent) WriteKey(buff core.DataBuffer) error {
	
	if err := buff.WriteString(s.ETag); err != nil {
		return err
	}
	if err := buff.WriteString(s.Name); err != nil {
		return err
	}
	return nil
}

func (s *StatEvent) ReadKey(buff core.DataBuffer) error {
	tg, err := buff.ReadString()
	if err != nil {
		return err
	}
	if tg != s.Tag() {
		return errors.New("etag not match")
	}
	tag, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.ETag = tag
	name, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.Name = name
	return nil
}

func (s *StatEvent) Read(buff core.DataBuffer) error {
	cnt, err := buff.ReadInt32()
	if err != nil {
		return err
	}
	s.Count = uint32(cnt)
	return nil

}

func (s *StatEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt32(int32(s.Count)); err != nil {
		return err
	}
	return nil
}




