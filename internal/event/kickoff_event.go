package event

import "gameclustering.com/internal/core"

type KickoffEvent struct {
	Source   string
	SystemId int64
	EventObj
}

func (s *KickoffEvent) ClassId() int {
	return KICKOFF_CID
}

func (s *KickoffEvent) ETag() string {
	return KICKOFF_ETAG
}

func (s *KickoffEvent) RecipientId() int64 {
	return s.SystemId
}

func (s *KickoffEvent) WriteKey(buff core.DataBuffer) error {
	if err := buff.WriteString(s.ETag()); err != nil {
		return err
	}
	return nil
}

func (s *KickoffEvent) ReadKey(buff core.DataBuffer) error {
	_, err := buff.ReadString()
	if err != nil {
		return err
	}
	return nil
}

func (s *KickoffEvent) Read(buff core.DataBuffer) error {
	source, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.Source = source
	sysId, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.SystemId = sysId
	return nil
}

func (s *KickoffEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteString(s.Source); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.SystemId); err != nil {
		return err
	}
	return nil
}

func (s *KickoffEvent) Outbound(buff core.DataBuffer) error {
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

func (s *KickoffEvent) Inbound(buff core.DataBuffer) error {
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
