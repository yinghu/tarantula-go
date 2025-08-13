package main

import (
	"encoding/json"

	"gameclustering.com/internal/event"
)

type Segment struct {
	InstanceId int64  `json:"-"`
	Name       string `json:"Name"`
}

type SegementSchedule struct {
	TournamentId int64  `json:"TournamentId,string"`
	Name         string `json:"Name"`
	Schedule
	Segments []Segment `json:"-"`

	*TournamentService `json:"-"`
}

func (t SegementSchedule) Join(join event.TournamentEvent) error {
	seg := t.Segments[0]
	join.InstanceId = seg.InstanceId
	t.updateSegment(join)
	return nil
}

func (t SegementSchedule) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}
