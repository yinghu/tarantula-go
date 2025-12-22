package main

import (
	"net/http"
	"strconv"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type TableInfo struct {
	Message string `json:"Message"`
	Code    int32  `json:"Code"`
	TableId int64  `json:"TableId"`
}

type MahjongTableSelector struct {
	*MahjongService
}

func (s *MahjongTableSelector) AccessControl() int32 {
	return bootstrap.PROTECTED_ACCESS_CONTROL
}

func (s *MahjongTableSelector) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	lobby := r.PathValue("lobby")
	sysId, err := strconv.ParseInt(r.PathValue("systemId"), 10, 64)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		ti := TableInfo{Message: err.Error(), Code: 400100, TableId: 0}
		w.Write(util.ToJson(ti))
		return
	}
	tc := make(chan TableInfo, 1)
	defer close(tc)
	pt := MahjongPlayToken{SystemId: sysId, Lobby: lobby, TableSelector: tc, Cmd: CMD_TABLE_PICK}
	s.Dispatcher <- pt
	ti := <-tc
	w.Write(util.ToJson(ti))
}
