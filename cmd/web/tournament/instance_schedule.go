package main

type InstanceSchedule struct {
	TournamentId      int64  `json:"-"`
	Name              string `json:"Name"`
	MaxEntries        int32  `json:"MaxEntries"`
	DurationInMinutes int32  `json:"DurationInMinutes"`
}
