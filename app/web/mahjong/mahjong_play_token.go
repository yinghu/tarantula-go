package main

import "gameclustering.com/internal/core"

const (
	CMD_SIT       int = 0
	CMD_DICE      int = 1
	CMD_DEAL      int = 2
	CMD_DRAW      int = 3
	CMD_DISCHARGE int = 4
	CMD_CLAIM     int = 5
	CMD_END       int = 99

	//internal
	CMD_JOINED int = 100
	CMD_LEFT   int = 101
)

type MahjongPlayToken struct {
	Table      int64
	SystemId   int64
	Cmd        int
	Seat       int
	Discharged int
}

func (mp *MahjongPlayToken) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(mp.Table); err != nil {
		return err
	}
	if err := buff.WriteInt64(mp.SystemId); err != nil {
		return err
	}
	if err := buff.WriteInt32(int32(mp.Cmd)); err != nil {
		return err
	}
	if err := buff.WriteInt32(int32(mp.Seat)); err != nil {
		return err
	}
	if err := buff.WriteInt32(int32(mp.Discharged)); err != nil {
		return err
	}
	return nil
}

func (mp *MahjongPlayToken) Read(buff core.DataBuffer) error {
	tbl, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	mp.Table = tbl
	sysId, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	mp.SystemId = sysId
	cmd, err := buff.ReadInt32()
	if err != nil {
		return err
	}
	mp.Cmd = int(cmd)
	seat, err := buff.ReadInt32()
	if err != nil {
		return err
	}
	mp.Seat = int(seat)
	discharged, err := buff.ReadInt32()
	if err != nil {
		return err
	}
	mp.Discharged = int(discharged)
	return nil
}
