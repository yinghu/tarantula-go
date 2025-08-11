package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type JoinType struct {
	Type string `json:"Type"`
}

type TournamentJoin struct {
	*TournamentService
}

func (s *TournamentJoin) AccessControl() int32 {
	return bootstrap.PROTECTED_ACCESS_CONTROL
}

func (s *TournamentJoin) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var join JoinType
	json.NewDecoder(r.Body).Decode(&join)
	w.WriteHeader(http.StatusOK)
	w.Write(util.ToJson(join))
}
