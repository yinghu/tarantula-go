package main

import "gameclustering.com/internal/event"

type InstanceSchedule struct {
	TournamentId      int64  `json:"-"`
	Name              string `json:"Name"`
	MaxEntries        int32  `json:"MaxEntries"`
	DurationInMinutes int32  `json:"DurationInMinutes"`
	Schedule

	*TournamentService
}

func (t InstanceSchedule) Join(join event.TournamentEvent) error {
	return nil
}
