package event

import (
	"time"

	"gameclustering.com/internal/core"
)

type NodeStartEvent struct {
	Id        int64     `json:"id,string"`
	NodeName  string    `json:"NodeName"`
	StartTime time.Time `json:"StartTime"`
	EventObj
}

func (s *NodeStartEvent) ClassId() int {
	return NODE_START_CID
}

func (s *NodeStartEvent) ETag() string {
	return NODE_START_ETAG
}

func (s *NodeStartEvent) WriteKey(buff core.DataBuffer) error {
	if err := buff.WriteString(s.ETag()); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.Id); err != nil {
		return err
	}
	return nil
}

func (s *NodeStartEvent) ReadKey(buff core.DataBuffer) error {
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

func (s *NodeStartEvent) Read(buff core.DataBuffer) error {
	nodeName, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.NodeName = nodeName
	tm, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.StartTime = time.UnixMilli(tm)
	return nil
}

func (s *NodeStartEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteString(s.NodeName); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.StartTime.UnixMilli()); err != nil {
		return err
	}

	return nil
}

func (s *NodeStartEvent) Outbound(buff core.DataBuffer) error {
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

func (s *NodeStartEvent) Inbound(buff core.DataBuffer) error {
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
