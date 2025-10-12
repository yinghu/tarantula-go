package event

import (
	"time"

	"gameclustering.com/internal/core"
)

type RegisterEvent struct {
	Name      string    `json:"login"`
	SystemId  int64     `json:"systemId,string"`
	RegisterTime time.Time `json:"registerTime"`
	EventObj
}

func (s *RegisterEvent) ClassId() int {
	return REGISTER_CID
}

func (s *RegisterEvent) ETag() string {
	return REGISTER_ETAG
}
func (s *RegisterEvent) Read(buffer core.DataBuffer) error {
	name, err := buffer.ReadString()
	if err != nil {
		return err
	}
	s.Name = name
	lt, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	s.RegisterTime = time.UnixMilli(lt)
	return nil
}

func (s *RegisterEvent) Write(buffer core.DataBuffer) error {
	if err := buffer.WriteString(s.Name); err != nil {
		return err
	}
	if err := buffer.WriteInt64(s.RegisterTime.UnixMilli()); err != nil {
		return err
	}
	return nil
}

func (s *RegisterEvent) ReadKey(buffer core.DataBuffer) error {
	_, err := buffer.ReadString()
	if err != nil {
		return err
	}
	sid, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	s.SystemId = sid
	id, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	s.OnOId(id)
	return nil
}

func (s *RegisterEvent) WriteKey(buffer core.DataBuffer) error {
	if err := buffer.WriteString(s.ETag()); err != nil {
		return err
	}
	if err := buffer.WriteInt64(s.SystemId); err != nil {
		return err
	}
	if err := buffer.WriteInt64(s.OId()); err != nil {
		return err
	}
	return nil
}

func (s *RegisterEvent) Outbound(buff core.DataBuffer) error {
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

func (s *RegisterEvent) Inbound(buff core.DataBuffer) error {
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
