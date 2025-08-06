package event

import (
	"gameclustering.com/internal/core"
)

type LoginEvent struct {
	Id            int32  `json:"-"`
	Name          string `json:"login"`
	Hash          string `json:"password"`
	ReferenceId   int32  `json:"referenceId"`
	SystemId      int64
	AccessControl int32      `json:"accessControl"`
	Cc            chan Chunk `json:"-"`
	EventObj      `json:"-"`
}

func (s *LoginEvent) ClassId() int {
	return LOGIN_CID
}


func (s *LoginEvent) Read(buffer core.DataBuffer) error {
	hash, err := buffer.ReadString()
	if err != nil {
		return err
	}
	s.Hash = hash
	refId, err := buffer.ReadInt32()
	if err != nil {
		return err
	}
	s.ReferenceId = refId
	sysId, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	s.SystemId = sysId
	return nil
}

func (s *LoginEvent) Write(buffer core.DataBuffer) error {
	err := buffer.WriteString(s.Hash)
	if err != nil {
		return err
	}
	err = buffer.WriteInt32(s.ReferenceId)
	if err != nil {
		return err
	}
	err = buffer.WriteInt64(s.SystemId)
	if err != nil {
		return err
	}
	return nil
}

func (s *LoginEvent) ReadKey(buffer core.DataBuffer) error {
	name, err := buffer.ReadString()
	if err != nil {
		return err
	}
	s.Name = name
	return nil
}

func (s *LoginEvent) WriteKey(buffer core.DataBuffer) error {
	err := buffer.WriteString(s.Name)
	if err != nil {
		return err
	}
	return nil
}

func (s *LoginEvent) Outbound(buff core.DataBuffer) error {
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

func (s *LoginEvent) Inbound(buff core.DataBuffer) error {
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
