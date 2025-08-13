package main

import "gameclustering.com/internal/event"

type Segment struct {
	InstanceId int64  `json:"-"`
	Name       string `json:"Name"`
}

type SegementSchedule struct {
	TournamentId int64  `json:"-"`
	Name         string `json:"Name"`
	Schedule
	Segments []Segment `json:"-"`

	*TournamentService
}

func (t SegementSchedule) Join(join event.TournamentEvent) error {
	seg := t.Segments[0]
	join.InstanceId = seg.InstanceId
	t.updateInstance(join)
	return nil
}
