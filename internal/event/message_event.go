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
	return 1
}

func (s *MessageEvent) Inbound(buff core.DataBuffer) {
	title, err := buff.ReadString()
	if err != nil {
		s.Listener().OnError(err)
		return
	}
	s.Title = title
	message, err := buff.ReadString()
	if err != nil {
		s.Listener().OnError(err)
		return
	}
	s.Message = message
	s.Listener().OnEvent(s)
}

func (s *MessageEvent) Outbound(buff core.DataBuffer) {
	buff.WriteString(s.Title)
	buff.WriteString(s.Message)
}
