package bootstrap

import "gameclustering.com/internal/core"

type ClusterObject struct {
	Key      string `json:"key"`
	Type     string `json:"type"`
	Balance  int32  `json:"balance"`
	SystemId int64  `json:"systemId,string"`
	Mutable  bool   `json:"mutable"`

	core.PersistentableObj
}

func (c *ClusterObject) FactoryId() uint32 {
	return 100
}

func (c *ClusterObject) ClassId() uint32 {
	return 1000
}

func (c *ClusterObject) WriteKey(buff core.DataBuffer) error {
	if err := buff.WriteString(c.Key); err != nil {
		return err
	}
	return nil
}

func (c *ClusterObject) ReadKey(buff core.DataBuffer) error {
	k, err := buff.ReadString()
	if err != nil {
		return err
	}
	c.Key = k
	return nil
}

func (s *ClusterObject) Read(buff core.DataBuffer) error {
	tp, err := buff.ReadString()
	if err != nil {
		return err
	}
	s.Type = tp
	balance, err := buff.ReadInt32()
	if err != nil {
		return err
	}
	s.Balance = balance
	sysId, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.SystemId = sysId
	return nil
}

func (s *ClusterObject) Write(buff core.DataBuffer) error {
	if err := buff.WriteString(s.Type); err != nil {
		return err
	}
	if err := buff.WriteInt32(s.Balance); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.SystemId); err != nil {
		return err
	}
	return nil
}
