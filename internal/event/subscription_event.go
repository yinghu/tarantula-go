package event

import (
	"gameclustering.com/internal/core"
)

type SubscriptionEvent struct {
	Id       int32  `json:"Id"`
	App      string `json:"App"`
	Name     string `json:"Name"`
	EventObj `json:"-"`
}

func (s *SubscriptionEvent) ClassId() int {
	return SUBSCRIPTION_CID
}

func (s *SubscriptionEvent) ETag() string {
	return SUBSCRIPTION_ETAG
}

func (s *SubscriptionEvent) Read(buff core.DataBuffer) error {
	id, err := buff.ReadInt32()
	if err != nil {
		return err
	}
	s.Id = id
	app, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.App = app
	name, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.Name = name
	return nil
}

func (s *SubscriptionEvent) Write(buff core.DataBuffer) error {
	err := buff.WriteInt32(s.Id)
	if err != nil {
		return err
	}
	err = buff.WriteString(s.App)
	if err != nil {
		return err
	}
	err = buff.WriteString(s.Name)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionEvent) Outbound(buff core.DataBuffer) error {
	err := s.Write(buff)
	if err != nil {
		s.Callback.OnError(s,err)
		return err
	}
	return nil
}

func (s *SubscriptionEvent) Inbound(buff core.DataBuffer) error {
	err := s.Read(buff)
	if err != nil {
		s.Callback.OnError(s,err)
		return err
	}
	s.Callback.OnEvent(s)
	return nil
}
