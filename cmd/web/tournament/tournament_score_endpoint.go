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

type TournamentScore struct {
	*TournamentService
}

func (s *TournamentScore) AccessControl() int32 {
	return bootstrap.PROTECTED_ACCESS_CONTROL
}

func (s *TournamentScore) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
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
	joined, err := tmnt.Score(join)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	join.InstanceId = joined.InstanceId
	w.Write(util.ToJson(join))
}
