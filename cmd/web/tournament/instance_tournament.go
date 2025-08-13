package main

import (
	"encoding/json"

	"gameclustering.com/internal/event"
)

type InstanceSchedule struct {
	TournamentId      int64  `json:"TournamentId,string"`
	Name              string `json:"Name"`
	MaxEntries        int32  `json:"MaxEntries"`
	DurationInMinutes int32  `json:"DurationInMinutes"`
	Schedule

	*TournamentService `json:"-"`
}

func (t InstanceSchedule) Join(join event.TournamentEvent) error {
	t.updateInstance(join, t.MaxEntries)
	return nil
}

func (t InstanceSchedule) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}
