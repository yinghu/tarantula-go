package event

import (
	"errors"

	"gameclustering.com/internal/core"
)

type IndexEvent struct {
	Id  int64  `json:"id,string"`
	Tag string `json:"Tag"`

	Index    []byte `json:"-"`
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
	if err := buff.WriteInt32(int32(len(s.Index))); err != nil {
		return err
	}
	if err := buff.Write(s.Index); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.Id); err != nil {
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
	sz, err := buff.ReadInt32()
	if err != nil {
		return err
	}
	index, err := buff.Read(int(sz))
	if err != nil {
		return err
	}
	s.Index = index
	id, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.Id = id
	return nil
}

func (s *IndexEvent) Read(buff core.DataBuffer) error {
	key, err := buff.Read(0)
	if err != nil {
		return err
	}
	s.Key = key
	return nil

}

func (s *IndexEvent) Write(buff core.DataBuffer) error {
	if err := buff.Write(s.Key); err != nil {
		return err
	}
	return nil
}

