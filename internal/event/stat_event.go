package event

import (
	"errors"

	"gameclustering.com/internal/core"
)

type StatEvent struct {
	Tag      string `json:"Tag"`
	Name     string `json:"Name"`
	Count    uint32 `json:"Count"`
	EventObj `json:"-"`
}

func (s *StatEvent) ClassId() int {
	return STAT_CID
}

func (s *StatEvent) ETag() string {
	return STAT_ETAG
}

func (s *StatEvent) WriteKey(buff core.DataBuffer) error {
	if err := buff.WriteString(s.ETag()); err != nil {
		return err
	}
	if err := buff.WriteString(s.Tag); err != nil {
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
	if tg != s.ETag() {
		return errors.New("etag not match")
	}
	tag, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.Tag = tag
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




