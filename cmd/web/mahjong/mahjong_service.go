package main 

import (
	"gameclustering.com/internal/bootstrap"
)
type MahjongService struct {
	bootstrap.AppManager
}

func (s *MahjongService) Config() string {
	return "/etc/tarantula/mahjong-conf.json"
}