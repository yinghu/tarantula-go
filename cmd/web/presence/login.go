package main

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

const (
	DB_OP_ERR_CODE int = 500100

	WRONG_PASS_CODE    int    = 400100
	WRONG_PASS_MSG     string = "wrong user/password"
	INVALID_TOKEN_CODE int    = 400101
	INVALID_TOKEN_MSG  string = "invalid token"
)

type OnSession struct {
	Successful bool   `json:"successful"`
	ErrorCode  int    `json:"errorCode"`
	Message    string `json:"message"`
	SystemId   int64  `json:"systemId"`
	Stub       int64  `json:"stub"`
	Token      string `json:"token"`
	Home       string `json:"home"`
}

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
	a, _ := buffer.ReadInt32()
	b, _ := buffer.ReadInt32()
	fmt.Printf("Reading from buffer %d, %d\n", a, b)
	return nil
}

func (s *Login) Write(buffer core.DataBuffer) error {
	buffer.WriteInt32(100)
	buffer.WriteInt32(200)
	return nil
}

func (s *Login) Inbound(buff core.DataBuffer) {
	s.Read(buff)
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

func errorMessage(msg string, code int) []byte {
	m := OnSession{Message: msg, ErrorCode: code}
	return util.ToJson(m)
}
func successMessage(msg string) []byte {
	m := OnSession{Message: msg, Successful: true}
	return util.ToJson(m)
}
