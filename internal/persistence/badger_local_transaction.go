package persistence

import (
	"fmt"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	badger "github.com/dgraph-io/badger/v4"
)

type BadgerLocalTransaction struct {
	ctx   *badger.Txn
	key   core.DataBuffer
	value core.DataBuffer
}

func (t *BadgerLocalTransaction) Get(p core.Persistentable) error {
	t.key.Clear()
	t.value.Clear()
	err := p.WriteKey(t.key)
	if err != nil {
		return err
	}
	t.key.Flip()
	pk, err := t.key.Read(0)
	if err != nil {
		return err
	}
	item, err := t.ctx.Get(pk)
	if err != nil {
		return err
	}
	err = item.Value(func(val []byte) error {
		err = t.value.Write(val)
		if err != nil {
			return err
		}
		t.value.Flip()
		return nil
	})
	if err != nil {
		return err
	}
	cid, err := t.value.ReadUInt32()
	if err != nil {
		return err
	}
	if p.ClassId() != cid {
		return fmt.Errorf("class id not matched %d , %d", cid, p.ClassId())
	}
	crv, err := t.value.ReadUInt64()
	if err != nil {
		return err
	}
	ctm, err := t.value.ReadUInt64()
	if err != nil {
		return err
	}
	p.Read(t.value)
	p.OnRevision(crv)
	p.OnTimestamp(ctm)
	return nil
}
func (t *BadgerLocalTransaction) Set(p core.Persistentable) error {
	t.key.Clear()
	t.value.Clear()
	if err := p.WriteKey(t.key); err != nil {
		return err
	}
	t.key.Flip()
	pk, err := t.key.Read(0)
	if err != nil {
		return err
	}
	item, err := t.ctx.Get(pk)
	var rev uint64 = 0
	if err == nil {
		err = item.Value(func(val []byte) error {
			return t.value.Write(val)
		})
		if err == nil {
			t.value.Flip()
			cid, err := t.value.ReadUInt32()
			if err != nil {
				return err
			}
			if cid != p.ClassId() {
				return fmt.Errorf("class id not matched %d , %d", cid, p.ClassId())
			}
			crv, err := t.value.ReadUInt64()
			if err != nil {
				return err
			}
			if crv != p.Revision() {
				return fmt.Errorf("class id not matched %d , %d", cid, p.ClassId())
			}
			rev = crv
		}
	}
	t.value.Clear()
	if err := t.value.WriteUInt32(p.ClassId()); err != nil {
		return err
	}
	if err := t.value.WriteUInt64(p.Revision() + 1); err != nil {
		return err
	}
	if err := t.value.WriteUInt64(p.Timestamp()); err != nil {
		return err
	}
	if err := p.Write(t.value); err != nil {
		return err
	}
	t.value.Flip()
	pv, err := t.value.Read(0)
	if err != nil {
		return err
	}

	t.ctx.Set(pk, pv)
	if rev > 0 {
		return nil
	}
	se := event.StatEvent{Tag: p.ETag(), Name: event.STAT_TOTAL}
	t.key.Clear()
	t.value.Clear()
	se.WriteKey(t.key)
	t.key.Flip()
	sk, _ := t.key.Read(0)
	si, err := t.ctx.Get(sk)
	if err == nil {
		err = si.Value(func(val []byte) error {
			err = t.value.Write(val)
			if err != nil {
				return err
			}
			t.value.Flip()
			t.value.ReadInt32()
			t.value.ReadInt64()
			t.value.ReadInt64()
			se.Read(t.value)
			se.Count = se.Count + 1
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		se.Count = 1
	}
	t.value.Clear()
	t.value.WriteUInt32(se.ClassId())
	t.value.WriteUInt64(se.Revision())
	t.value.WriteUInt64(uint64(time.Now().UnixMilli()))
	se.Write(t.value)
	t.value.Flip()
	sv, _ := t.value.Read(0)
	return t.ctx.Set(sk, sv)
}
func (t *BadgerLocalTransaction) Del(p core.Persistentable) error {
	t.key.Clear()
	t.value.Clear()
	err := p.ReadKey(t.key)
	if err != nil {
		return err
	}
	t.key.Flip()
	pk, err := t.key.Read(0)
	if err != nil {
		return err
	}
	return t.ctx.Delete(pk)
}

func (t *BadgerLocalTransaction) Commit() error {
	return t.ctx.Commit()
}

func (t *BadgerLocalTransaction) Rollback() {
	t.ctx.Discard()
}
