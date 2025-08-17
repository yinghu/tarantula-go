package event

import (
	"time"

	"gameclustering.com/internal/core"
)

type LoginEvent struct {
	Id        int64     `json:"id"`
	Name      string    `json:"login"`
	SystemId  int64     `json:"systemId:string"`
	LoginTime time.Time `json:"loginTime"`
	EventObj
}

func (s *LoginEvent) ClassId() int {
	return LOGIN_CID
}

func (s *LoginEvent) ETag() string {
	return LOGIN_ETAG
}
func (s *LoginEvent) Read(buffer core.DataBuffer) error {
	name, err := buffer.ReadString()
	if err != nil {
		return err
	}
	s.Name = name
	lt, err := buffer.ReadInt64()
	if err!=nil{
		return err
	}
	s.LoginTime = time.UnixMilli(lt)
	return nil
}

func (s *LoginEvent) Write(buffer core.DataBuffer) error {
	if err := buffer.WriteString(s.Name); err != nil {
		return err
	}
	if err := buffer.WriteInt64(s.LoginTime.UnixMilli()); err != nil {
		return err
	}
	return nil
}

func (s *LoginEvent) ReadKey(buffer core.DataBuffer) error {
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
	s.Id = id
	return nil
}

func (s *LoginEvent) WriteKey(buffer core.DataBuffer) error {
	if err := buffer.WriteString(s.ETag()); err != nil {
		return err
	}
	if err := buffer.WriteInt64(s.SystemId); err != nil {
		return err
	}
	if err := buffer.WriteInt64(s.Id); err != nil {
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
