package main

import "gameclustering.com/internal/core"

type KeyIndex struct {
	Prefix uint32      `json:"prefix"`
	Key    []byte      `json:"-"`
	Master core.Node   `json:"master"`
	Slaves []core.Node `json:"slaves"`

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
	if err := buffer.WriteUInt32(k.Master.RingToken); err != nil {
		return err
	}
	if err := buffer.WriteInt32(int32(len(k.Slaves))); err != nil {
		return err
	}
	for _, n := range k.Slaves {
		if err := buffer.WriteUInt32(n.RingToken); err != nil {
			return err
		}
	}
	return nil
}
func (k *KeyIndex) Read(buffer core.DataBuffer) error {
	master, err := buffer.ReadUInt32()
	if err != nil {
		return err
	}
	k.Master = core.Node{RingToken: master}
	slen, err := buffer.ReadInt32()
	if err != nil {
		return err
	}
	k.Slaves = make([]core.Node, 0, slen)
	for {
		if slen == 0 {
			break
		}
		slave, err := buffer.ReadUInt32()
		if err != nil {
			return err
		}
		k.Slaves = append(k.Slaves, core.Node{RingToken: slave})
		slen--
	}
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
