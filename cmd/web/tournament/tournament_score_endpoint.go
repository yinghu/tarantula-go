package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type TournamentScore struct {
	*TournamentService
}

func (s *TournamentScore) AccessControl() int32 {
	return bootstrap.PROTECTED_ACCESS_CONTROL
}

func (s *TournamentScore) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var score event.TournamentEvent
	json.NewDecoder(r.Body).Decode(&score)
	w.WriteHeader(http.StatusOK)
	w.Write(util.ToJson(score))
}
