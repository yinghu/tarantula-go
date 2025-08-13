package main

import (
	"encoding/json"
	"fmt"

	"gameclustering.com/internal/event"
)

type InstanceSchedule struct {
	TournamentId      int64  `json:"TournamentId,string"`
	Name              string `json:"Name"`
	MaxEntries        int32  `json:"MaxEntries"`
	DurationInMinutes int32  `json:"DurationInMinutes"`
	TotalJoined       int32  `json:"TotalJoined"`
	Schedule
	*TournamentService `json:"-"`
}

func (t *InstanceSchedule) Start() error {
	return nil
}

func (t *InstanceSchedule) Score(join event.TournamentEvent) (event.TournamentEvent, error) {
	return join,nil
}

func (t *InstanceSchedule) Join(join event.TournamentEvent) (event.TournamentEvent, error) {
	total, err := t.updateInstance(join, t.MaxEntries)
	if err != nil {
		return join, err
	}
	t.TotalJoined += total
	return join, nil
}

func (t *InstanceSchedule) MarshalJSON() ([]byte, error) {
	data := make(map[string]any)
	data["TournamentId"] = fmt.Sprintf("%d", t.TournamentId)
	data["Name"] = t.Name
	data["MaxEntries"] = t.MaxEntries
	data["DurationInMinutes"] = t.DurationInMinutes
	data["TotalJoined"] = t.TotalJoined
	data["Schedule"] = t.Schedule
	return json.Marshal(data)
}
