package event

import (
	"time"

	"gameclustering.com/internal/core"
)

type MessageEvent struct {
	Title    string    `json:"title"`
	Message  string    `json:"message"`
	DateTime time.Time `json:"dataTime"`
	Id       int64     `json:"id,string"`
	Source   string    `json:"source"`
	EventObj
}

func (s *MessageEvent) ClassId() int {
	return MESSAGE_CID
}

func (s *MessageEvent) ETag() string {
	return MESSAGE_ETAG
}

func (s *MessageEvent) WriteKey(buff core.DataBuffer) error {
	if err := buff.WriteString(s.ETag()); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.Id); err != nil {
		return err
	}
	return nil
}

func (s *MessageEvent) ReadKey(buff core.DataBuffer) error {
	_, err := buff.ReadString()
	if err != nil {
		return err
	}
	id, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.Id = id
	return nil
}

func (s *MessageEvent) Read(buff core.DataBuffer) error {
	title, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.Title = title
	message, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.Message = message
	tm, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.DateTime = time.UnixMilli(tm)
	source, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.Source = source
	return nil
}

func (s *MessageEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteString(s.Title); err != nil {
		return err
	}
	if err := buff.WriteString(s.Message); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.DateTime.UnixMilli()); err != nil {
		return err
	}
	if err := buff.WriteString(s.Source); err != nil {
		return err
	}
	return nil
}

func (s *MessageEvent) Outbound(buff core.DataBuffer) error {
	err := s.WriteKey(buff)
	if err != nil {
		s.Callback.OnError(s,err)
		return err
	}
	err = s.Write(buff)
	if err != nil {
		s.Callback.OnError(s,err)
		return err
	}
	return nil
}

func (s *MessageEvent) Inbound(buff core.DataBuffer) error {
	err := s.ReadKey(buff)
	if err != nil {
		s.Callback.OnError(s,err)
		return err
	}
	err = s.Read(buff)
	if err != nil {
		s.Callback.OnError(s,err)
		return err
	}
	s.Callback.OnEvent(s)
	return nil
}
