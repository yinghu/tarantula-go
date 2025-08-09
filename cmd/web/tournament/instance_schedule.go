package main

import "time"

type Schedule struct {
	StartTime time.Time `json:"StartTime"`
	CloseTime time.Time `json:"CloseTime"`
	EndTime   time.Time `json:"EndTime"`
}

type InstanceSchedule struct {
	TournamentId      int64  `json:"-"`
	Name              string `json:"Name"`
	MaxEntries        int32  `json:"MaxEntries"`
	DurationInMinutes int32  `json:"DurationInMinutes"`
	Schedule
}
