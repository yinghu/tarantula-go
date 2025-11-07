package main

import "gameclustering.com/internal/mj"

type MahjongTable struct {
	Setup   mj.ClassicMahjong
	Players []MahjongPlayer
}

func (m *MahjongTable) Sit(systemId int64, sn int) error {
	return nil
}

func (m *MahjongTable) Dice() error {
	return nil
}
func (m *MahjongTable) Deal() error {
	return nil
}

func (m *MahjongTable) Discharge() error {
	return nil
}

func (m *MahjongTable) Chow() error {
	return nil
}
