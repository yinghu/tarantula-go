package event

import (
	"net"

	"gameclustering.com/internal/core"
)

type JoinEvent struct {
	Ticket   string `json:"Ticket"`
	Client   net.Conn
	SystemId int64
	Flag     int64
	core.EventObj
}

func (s *JoinEvent) ClassId() uint32 {
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
	ticket, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.Ticket = ticket
	flag, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.Flag = flag
	return nil
}

func (s *JoinEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteString(s.Ticket); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.Flag); err != nil {
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
	//s.Callback.OnEvent(s)
	return nil
}

func (s *JoinEvent) RecipientId() int64 {
	return s.SystemId
}
