package event

import (
	"net"

	"gameclustering.com/internal/core"
)

type JoinEvent struct {
	Token    string `json:"Token"`
	Client   net.Conn
	Pending  core.DataBuffer
	SystemId int64
	EventObj
}

func (s *JoinEvent) ClassId() int {
	return JOIN_CID
}

func (s *JoinEvent) ETag() string {
	return JOIN_ETAG
}

func (s *JoinEvent) WriteKey(buff core.DataBuffer) error {
	if err := buff.WriteString(s.ETag()); err != nil {
		return err
	}
	return nil
}

func (s *JoinEvent) ReadKey(buff core.DataBuffer) error {
	_, err := buff.ReadString()
	if err != nil {
		return err
	}
	return nil
}

func (s *JoinEvent) Read(buff core.DataBuffer) error {
	token, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.Token = token
	return nil
}

func (s *JoinEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteString(s.Token); err != nil {
		return err
	}
	return nil
}

func (s *JoinEvent) Outbound(buff core.DataBuffer) error {
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

func (s *JoinEvent) Inbound(buff core.DataBuffer) error {
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
