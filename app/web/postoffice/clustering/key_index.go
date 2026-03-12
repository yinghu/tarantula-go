package clustering

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

type KeyIndex struct {
	Prefix uint32           `json:"prefix"`
	Key    []byte           `json:"-"`
	Header *protocol.Header `json:"header"`
	Nodes  []uint32         `json:"-"`

	core.PersistentableObj `json:"-"`
}

func (k *KeyIndex) WriteKey(buffer core.DataBuffer) error {
	if err := buffer.WriteUInt32(k.Prefix); err != nil {
		return err
	}
	if err := buffer.WriteUInt32(k.Header.FactoryId); err != nil {
		return err
	}
	if err := buffer.WriteUInt32(k.Header.ClassId); err != nil {
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
	if err := buffer.WriteUInt64(k.Header.Revision); err != nil {
		return err
	}
	if err := buffer.WriteUInt64(k.Header.Timestamp); err != nil {
		return err
	}
	if err := buffer.WriteUInt32(k.Header.Size); err != nil {
		return err
	}
	if err := buffer.WriteBool(k.Header.Mutable); err != nil {
		return err
	}
	return nil
}
func (k *KeyIndex) Read(buffer core.DataBuffer) error {
	rev, err := buffer.ReadUInt64()
	if err != nil {
		return err
	}
	k.Header.Revision = rev
	ts, err := buffer.ReadUInt64()
	if err != nil {
		return err
	}
	k.Header.Timestamp = ts
	sz, err := buffer.ReadUInt32()
	if err != nil {
		return err
	}
	k.Header.Size = sz
	mtu, err := buffer.ReadBool()
	if err != nil {
		return err
	}
	k.Header.Mutable = mtu
	return nil
}
func (k *KeyIndex) ReadKey(buffer core.DataBuffer) error {
	prefix, err := buffer.ReadUInt32()
	if err != nil {
		return err
	}
	k.Prefix = prefix
	fid, err := buffer.ReadUInt32()
	if err != nil {
		return err
	}
	k.Header.FactoryId = fid
	cid, err := buffer.ReadUInt32()
	if err != nil {
		return err
	}
	k.Header.ClassId = cid
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

//HELP METHODS

func (k *KeyIndex) Val(val []byte) error {
	kBuff := core.NewBuffer(COMPOSIT_KEY_MAX)
	if err := kBuff.Write(val); err != nil {
		return err
	}
	kBuff.Flip()
	return k.Read(kBuff)
}

func (k *KeyIndex) CompositKey() ([]byte, error) {
	kBuff := core.NewBuffer(COMPOSIT_KEY_MAX)
	if err := k.WriteKey(kBuff); err != nil {
		return nil, err
	}
	kBuff.Flip()
	return kBuff.Read(0)
}
func (k *KeyIndex) Value() ([]byte, error) {
	kBuff := core.NewBuffer(COMPOSIT_KEY_MAX)
	if err := k.Write(kBuff); err != nil {
		return nil, err
	}
	kBuff.Flip()
	return kBuff.Read(0)
}
func (k *KeyIndex) Pair() ([]byte, []byte, error) {
	kBuff := core.NewBuffer(COMPOSIT_KEY_MAX)
	if err := k.WriteKey(kBuff); err != nil {
		return nil, nil, err
	}
	kBuff.Flip()
	key, err := kBuff.Read(0)
	if err != nil {
		return nil, nil, err
	}
	kBuff.Clear()
	if err := k.Write(kBuff); err != nil {
		return nil, nil, err
	}
	kBuff.Flip()
	value, err := kBuff.Read(0)
	if err != nil {
		return nil, nil, err
	}
	return key, value, nil
}
