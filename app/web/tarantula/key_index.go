package main

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/tarantula/protocol"
)

const (
	COMPOSIT_KEY_MAX int = 500
)

type KeyIndex struct {
	Prefix                 uint32           `json:"prefix"`
	Key                    []byte           `json:"-"`
	Header                 *protocol.Header `json:"header"`
	core.PersistentableObj `json:"-"`
}

func (k *KeyIndex) WriteKey(buffer core.DataBuffer) error {
	if err := buffer.WriteUInt32(k.Prefix); err != nil {
		return err
	}
	if err := buffer.WriteInt32(int32(len(k.Key))); err != nil {
		return err
	}
	if err := buffer.Write(k.Key); err != nil {
		return err
	}
	return nil
}
func (k *KeyIndex) Write(buffer core.DataBuffer) error {
	if err := buffer.WriteInt64(k.Header.Revision); err != nil {
		return err
	}
	if err := buffer.WriteInt32(k.Header.FactoryId); err != nil {
		return err
	}
	if err := buffer.WriteInt32(k.Header.ClassId); err != nil {
		return err
	}
	if err := buffer.WriteInt64(k.Header.Timestamp); err != nil {
		return err
	}
	if err := buffer.WriteInt32(k.Header.Size); err != nil {
		return err
	}
	return nil
}
func (k *KeyIndex) Read(buffer core.DataBuffer) error {
	rev, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	k.Header.Revision = rev
	fid, err := buffer.ReadInt32()
	if err != nil {
		return err
	}
	k.Header.FactoryId = fid
	cid, err := buffer.ReadInt32()
	if err != nil {
		return err
	}
	k.Header.ClassId = cid
	ts, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	k.Header.Timestamp = ts
	sz, err := buffer.ReadInt32()
	if err != nil {
		return err
	}
	k.Header.Size = sz
	return nil
}
func (k *KeyIndex) ReadKey(buffer core.DataBuffer) error {
	prefix, err := buffer.ReadUInt32()
	if err != nil {
		return err
	}
	k.Prefix = prefix
	klen, err := buffer.ReadInt32()
	if err != nil {
		return err
	}
	key, err := buffer.Read(int(klen))
	if err != nil {
		return err
	}
	k.Key = key
	return nil
}

func (k *KeyIndex) CompositKey() ([]byte, error) {
	ksz := len(k.Key)
	if ksz+20 > COMPOSIT_KEY_MAX{
		return []byte{},fmt.Errorf("Key size overflow %d",ksz)
	}
	buffer := core.NewBuffer(COMPOSIT_KEY_MAX)
	buffer.WriteInt32(k.Header.FactoryId)
	buffer.WriteInt32(k.Header.ClassId)
	buffer.WriteInt32(int32(ksz))
	buffer.Write(k.Key)
	buffer.WriteInt64(k.Header.Revision)
	buffer.Flip()
	return buffer.Read(0)
}
