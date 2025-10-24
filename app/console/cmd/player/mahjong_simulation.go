package player

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type MahjongEvent struct {
	Cmd      string
	SystemId int64
	event.EventObj
}

func (s *MahjongEvent) ClassId() int {
	return 100
}

func (s *MahjongEvent) ETag() string {
	return "mj"
}

func (s *MahjongEvent) WriteKey(buff core.DataBuffer) error {
	if err := buff.WriteString(s.ETag()); err != nil {
		return err
	}
	return nil
}

func (s *MahjongEvent) ReadKey(buff core.DataBuffer) error {
	_, err := buff.ReadString()
	if err != nil {
		return err
	}
	return nil
}

func (s *MahjongEvent) Read(buff core.DataBuffer) error {
	sysId, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.SystemId = sysId
	cmd, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.Cmd = cmd
	return nil
}

func (s *MahjongEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(s.SystemId); err != nil {
		return err
	}
	if err := buff.WriteString(s.Cmd); err != nil {
		return err
	}
	return nil
}

func (s *MahjongEvent) Outbound(buff core.DataBuffer) error {
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

func (s *MahjongEvent) Inbound(buff core.DataBuffer) error {
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

func (s *MahjongEvent) RecipientId() int64 {
	return s.SystemId
}

type SampleCreator struct {
}

func (s *SampleCreator) Create(cid int, topic string) (event.Event, error) {
	e := event.CreateEvent(cid)
	if e != nil {
		e.OnTopic(topic)
		e.OnListener(s)
		return e, nil
	}
	me := MahjongEvent{}
	me.Callback = s
	return &me, nil
}

func (s *SampleCreator) OnError(e event.Event, err error) {
	fmt.Printf("On event error %v %s\n", e, err.Error())
}

func (s *SampleCreator) OnEvent(e event.Event) {
	fmt.Printf("On event %v\n", e)

}
