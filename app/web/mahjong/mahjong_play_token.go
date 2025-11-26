package main

import "gameclustering.com/internal/core"

const (
	CMD_SIT       int = 0
	CMD_DICE      int = 1
	CMD_DEAL      int = 2
	CMD_DRAW      int = 3
	CMD_DISCHARGE int = 4
	CMD_PONG      int = 5
	CMD_KONG      int = 6
	CMD_CHOW      int = 7
	CMD_CLAIM     int = 8
	CMD_TABLE     int = 9

	CMD_TURN_CONTINUE int = 96
	CMD_TURN_END      int = 97
	CMD_RESET         int = 98
	CMD_END           int = 99

	//internal

	CMD_JOINED int = 100
	CMD_LEFT   int = 101
)

type MahjongPlayToken struct {
	Id       int64
	TableId  int64
	SystemId int64
	Cmd      int
	Seat     int
	Selected int
}

func (mp *MahjongPlayToken) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(mp.Id); err != nil {
		return err
	}
	if err := buff.WriteInt64(mp.TableId); err != nil {
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
	if err := buff.WriteInt32(int32(mp.Selected)); err != nil {
		return err
	}
	return nil
}

func (mp *MahjongPlayToken) Read(buff core.DataBuffer) error {
	id, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	mp.Id = id
	tbl, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	mp.TableId = tbl
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
	selected, err := buff.ReadInt32()
	if err != nil {
		return err
	}
	mp.Selected = int(selected)
	return nil
}
