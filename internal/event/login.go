package event

import (
	"fmt"

	"gameclustering.com/internal/core"
)

type Login struct {
	Name          string `json:"login"`
	Hash          string `json:"password"`
	ReferenceId   int32  `json:"referenceId"`
	SystemId      int64
	AccessControl int32 `json:"accessControl"`
	EventObj            //Event default
}

func (s *Login) ClassId() int {
	return 10
}

func (s *Login) Read(buffer core.DataBuffer) error {
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

func (s *Login) Write(buffer core.DataBuffer) error {
	buffer.WriteString(s.Hash)
	buffer.WriteInt32(s.ReferenceId)
	buffer.WriteInt64(s.SystemId)
	return nil
}

func (s *Login) ReadKey(buffer core.DataBuffer) error {
	name, err := buffer.ReadString()
	if err != nil {
		return err
	}
	s.Name = name
	return nil
}

func (s *Login) WriteKey(buffer core.DataBuffer) error {
	buffer.WriteString(s.Name)
	return nil
}

func (s *Login) Inbound(buff core.DataBuffer) {
	s.ReadKey(buff)
	s.Read(buff)
	for {
		sz, err := buff.ReadInt32()
		if err != nil {
			s.streaming(Chunk{Remaining: true, Data: []byte{0}})
			break
		}
		if sz == 0 {
			s.streaming(Chunk{Remaining: true, Data: []byte{0}})
			break
		}

		pd, err := buff.Read(int(sz))
		if err != nil {
			s.streaming(Chunk{Remaining: true, Data: []byte{0}})
			break
		}
		s.streaming(Chunk{Remaining: true, Data: pd})
	}
	buff.WriteInt32(100)
	buff.WriteString("bye")
	s.Listener().OnEvent(s)
}

func (s *Login) Outbound(buff core.DataBuffer) {
	s.WriteKey(buff)
	s.Write(buff)
	buff.WriteInt32(12)
	buff.Write([]byte("login passed"))
	buff.WriteInt32(12)
	buff.Write([]byte("login passed"))
	buff.WriteInt32(0)
	buff.ReadInt32()
	buff.ReadString()

}

func (s *Login) streaming(c Chunk) {
}

func (s *Login) OnError(err error) {
	fmt.Printf("On error %s\n", err.Error())
}
