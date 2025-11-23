package main

import "gameclustering.com/internal/core"

type MahjongPlayTurn struct {
	Cmd       int
	CountDown int
}

func (mp *MahjongPlayTurn) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt32(int32(mp.Cmd)); err != nil {
		return err
	}
	if err := buff.WriteInt32(int32(mp.CountDown)); err != nil {
		return err
	}
	return nil
}
