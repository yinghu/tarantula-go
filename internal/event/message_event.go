package event

import (
	"gameclustering.com/internal/core"
)

type MessageEvent struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Id       int64  `json:"id,string"`
	EventObj `json:"-"`
}

func (s *MessageEvent) ClassId() int {
	return MESSAGE_CID
}

func (s *MessageEvent) ETag() string {
	return MESSAGE_ETAG
}

func (s *MessageEvent) WriteKey(buff core.DataBuffer) error {
	buff.WriteString(s.ETag())
	buff.WriteInt64(s.Id)
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
	return nil
}

func (s *MessageEvent) Write(buff core.DataBuffer) error {
	err := buff.WriteString(s.Title)
	if err != nil {
		return err
	}
	err = buff.WriteString(s.Message)
	if err != nil {
		return err
	}
	return nil
}

func (s *MessageEvent) Outbound(buff core.DataBuffer) error {
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

func (s *MessageEvent) Inbound(buff core.DataBuffer) error {
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
