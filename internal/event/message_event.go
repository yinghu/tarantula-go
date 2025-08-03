package event

import (
	"gameclustering.com/internal/core"
)

type MessageEvent struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	EventObj `json:"-"`
}

func (s *MessageEvent) ClassId() int {
	return MESSAGE_CID
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
	err := s.Write(buff)
	if err != nil {
		s.Callback.OnError(err)
		return err
	}
	s.Callback.OnEvent(s)
	return nil
}

func (s *MessageEvent) Inbound(buff core.DataBuffer) error {
	err := s.Read(buff)
	if err != nil {
		s.Callback.OnError(err)
		return err
	}
	s.Callback.OnEvent(s)
	return nil
}
