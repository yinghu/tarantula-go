package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type MahjongTableSelector struct {
	*MahjongService
}

func (s *MahjongTableSelector) AccessControl() int32 {
	return bootstrap.PROTECTED_ACCESS_CONTROL
}

func (s *MahjongTableSelector) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pts := s.Table.Setup.Dice()
	w.Write(util.ToJson(pts))
}
