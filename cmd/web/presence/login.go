package main

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type Login struct {
	Name        string `json:"login"`
	Hash        string `json:"password"`
	ReferenceId int32  `json:"referenceId"`
	SystemId    int64

	event.EventObj //Event default
}

func (s *Login) ClassId() int {
	return 10
}

func (s *Login) Read(buffer core.DataBuffer) error {
	name, err := buffer.ReadString()
	if err != nil {
		return err
	}
	s.Name = name
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
	buffer.WriteString(s.Name)
	buffer.WriteString(s.Hash)
	buffer.WriteInt32(s.ReferenceId)
	buffer.WriteInt64(s.SystemId)
	return nil
}

func (s *Login) Inbound(buff core.DataBuffer) {
	s.Read(buff)
	fmt.Printf("Login : %v\n", s)
	for {
		sz, err := buff.ReadInt32()
		if err != nil {
			s.streaming(event.Chunk{Remaining: true, Data: []byte{0}})
			break
		}
		if sz == 0 {
			s.streaming(event.Chunk{Remaining: true, Data: []byte{0}})
			break
		}
		fmt.Printf("Pending data %d\n", sz)
		pd, err := buff.Read(int(sz))
		if err != nil {
			fmt.Printf("Read err %s\n", err.Error())
			s.streaming(event.Chunk{Remaining: true, Data: []byte{0}})
			break
		}
		s.streaming(event.Chunk{Remaining: true, Data: pd})
	}
	buff.WriteInt32(100)
	buff.WriteString("bye")
}

func (s *Login) Outbound(buff core.DataBuffer) {
	s.Write(buff)
	buff.WriteInt32(12)
	buff.Write([]byte("login passed"))
	buff.WriteInt32(12)
	buff.Write([]byte("login passed"))
	buff.WriteInt32(0)
	r, _ := buff.ReadInt32()
	x, _ := buff.ReadString()
	fmt.Printf("%d %s\n", r, x)
}

func (s *Login) OnEvent(buff core.DataBuffer) {
	r, _ := buff.ReadInt32()
	x, _ := buff.ReadString()
	fmt.Printf("%d %s\n", r, x)
}

func (s *Login) streaming(c event.Chunk) {
	fmt.Printf("REV : %s\n", string(c.Data))
	//s.listener <- c
}

func (s *Login) OnError(err error) {
	fmt.Printf("On error %s\n", err.Error())
}
