package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type TournamentRaceBoard struct {
	*TournamentService
}

func (s *TournamentRaceBoard) AccessControl() int32 {
	return bootstrap.PROTECTED_ACCESS_CONTROL
}

func (s *TournamentRaceBoard) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var join event.TournamentEvent
	json.NewDecoder(r.Body).Decode(&join)
	tmnt := s.tournaments[join.TournamentId]
	w.WriteHeader(http.StatusOK)
	if tmnt == nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: fmt.Sprintf("tournament not available :%d", join.TournamentId)}))
		return
	}
	join.SystemId = rs.SystemId
	w.Write(util.ToJson(tmnt.Listing(join)))
}
