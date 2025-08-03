package event

import (
	"gameclustering.com/internal/core"
)

type Login struct {
	Id            int32  `json:"-"`
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





