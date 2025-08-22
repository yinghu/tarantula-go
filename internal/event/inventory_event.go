package event

import (
	"time"

	"gameclustering.com/internal/core"
)

type InventoryEvent struct {
	SystemId    int64     `json:"SystemId,string"`
	InventoryId int32     `json:"InventoryId"`
	ItemId      int64     `json:"ItemId,string"`
	Source      string    `json:"Source"`
	GrantTime   time.Time `json:"GrantTime"`
	Description string    `json:"Description"`
	EventObj    `json:"-"`
}

func (s *InventoryEvent) ClassId() int {
	return INVENTORY_CID
}

func (s *InventoryEvent) ETag() string {
	return INVENTORY_ETAG
}

func (s *InventoryEvent) WriteKey(buff core.DataBuffer) error {
	if err := buff.WriteString(s.ETag()); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.SystemId); err != nil {
		return err
	}
	if err := buff.WriteInt32(s.InventoryId); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.ItemId); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.OId()); err != nil {
		return err
	}
	return nil
}

func (s *InventoryEvent) ReadKey(buff core.DataBuffer) error {
	_, err := buff.ReadString()
	if err != nil {
		return err
	}
	systemId, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.SystemId = systemId
	inventoryId, err := buff.ReadInt32()
	if err != nil {
		return err
	}
	s.InventoryId = inventoryId
	itemId, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.ItemId = itemId
	id, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.OnOId(id)
	return nil
}

func (s *InventoryEvent) Read(buff core.DataBuffer) error {
	source, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.Source = source
	desc, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.Description = desc

	lt, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.GrantTime = time.UnixMilli(lt)
	return nil
}

func (s *InventoryEvent) Write(buff core.DataBuffer) error {
	if err := buff.WriteString(s.Source); err != nil {
		return err
	}
	if err := buff.WriteString(s.Description); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.GrantTime.UnixMilli()); err != nil {
		return err
	}
	return nil
}

func (s *InventoryEvent) Outbound(buff core.DataBuffer) error {
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

func (s *InventoryEvent) Inbound(buff core.DataBuffer) error {
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
