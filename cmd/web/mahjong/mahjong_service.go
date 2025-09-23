package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
)

type MahjongService struct {
	bootstrap.AppManager
	ClassicMahjong
}

func (s *MahjongService) Config() string {
	return "/etc/tarantula/mahjong-conf.json"
}

func (s *MahjongService) Start(f conf.Env, c core.Cluster) error {
	s.ItemUpdater = s
	s.AppManager.Start(f, c)
	s.ClassicMahjong = ClassicMahjong{}
	s.ClassicMahjong.New()
	http.Handle("/mahjong/dice", bootstrap.Logging(&MahjongDicer{MahjongService: s}))
	return nil
}
