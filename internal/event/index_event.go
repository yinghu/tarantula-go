package event

import (
	"errors"

	"gameclustering.com/internal/core"
)

type IndexEvent struct {
	Tag      string `json:"Tag"`
	
	Key      []byte `json:"-"`
	EventObj `json:"-"`
}

func (s *IndexEvent) ClassId() int {
	return INDEX_CID
}

func (s *IndexEvent) ETag() string {
	return INDEX_ETAG
}

func (s *IndexEvent) WriteKey(buff core.DataBuffer) error {
	if err := buff.WriteString(s.ETag()); err != nil {
		return err
	}
	if err := buff.WriteString(s.Tag); err != nil {
		return err
	}
	
	return nil
}

func (s *IndexEvent) ReadKey(buff core.DataBuffer) error {
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
	return nil
}

func (s *IndexEvent) Read(buff core.DataBuffer) error {
		
	return nil

}

func (s *IndexEvent) Write(buff core.DataBuffer) error {
	
	return nil
}

func (s *IndexEvent) Outbound(buff core.DataBuffer) error {
	err := s.WriteKey(buff)
	if err != nil {
		s.Callback.OnError(err)
		return err
	}
	err = s.Write(buff)
	if err != nil {
		s.Callback.OnError(err)
		return err
	}
	return nil
}

func (s *IndexEvent) Inbound(buff core.DataBuffer) error {
	err := s.ReadKey(buff)
	if err != nil {
		s.Callback.OnError(err)
		return err
	}
	err = s.Read(buff)
	if err != nil {
		s.Callback.OnError(err)
		return err
	}
	s.Callback.OnEvent(s)
	return nil
}
