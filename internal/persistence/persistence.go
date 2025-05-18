package persistence

import (
	"gameclustering.com/internal/core"
)

type PersistentableFactory interface {
	Create(classId int) Persistentable
}

type Persistentable interface {
	Write(value core.DataBuffer) error
	WriteKey(key core.DataBuffer) error
	Read(value core.DataBuffer) error
	ReadKey(key core.DataBuffer) error
	ClassId() int
	Revision() int64
}

type PersistentableObj struct {
	Fid int
	Cid int
	Rev int64
}

func (s *PersistentableObj) Write(value core.DataBuffer) error {
	return nil
}

func (s *PersistentableObj) WriteKey(value core.DataBuffer) error {
	return nil
}

func (s *PersistentableObj) Read(value core.DataBuffer) error {
	return nil
}

func (s *PersistentableObj) ReadKey(value core.DataBuffer) error {
	return nil
}

func (s *PersistentableObj) ClassId() int {
	return s.Cid
}

func (s *PersistentableObj) Revision() int64 {
	return s.Rev
}
