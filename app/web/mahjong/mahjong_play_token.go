package main

import "gameclustering.com/internal/core"

const (
	CMD_SIT  int = 0
	CMD_DICE int = 1
	CMD_PLAY int = 2
	CMD_END  int = 3
)

type MahjongPlayToken struct {
	SystemId int64
	Cmd      int
}

func (mp *MahjongPlayToken) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(mp.SystemId); err != nil {
		return err
	}
	if err := buff.WriteInt32(int32(mp.Cmd)); err != nil {
		return err
	}
	return nil
}

func (mp *MahjongPlayToken) Read(buff core.DataBuffer) error {
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
	return nil
}
