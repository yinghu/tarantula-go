package main

import (
	"encoding/json"
	"fmt"

	"gameclustering.com/internal/event"
)

type Segment struct {
	InstanceId  int64  `json:"-"`
	Name        string `json:"Name"`
	TotalJoined int32  `json:"TotalJoined"`
}

type SegmentSchedule struct {
	TournamentId int64  `json:"TournamentId,string"`
	Name         string `json:"Name"`
	TotalJoined  int32  `json:"TotalJoined"`
	Schedule
	Segments []*Segment `json:"Segments"`

	*TournamentService `json:"-"`
}

func (t *SegmentSchedule) Start() error {
	return nil
}

func (t *SegmentSchedule) Score(join event.TournamentEvent) (event.TournamentEvent, error) {
	return join, nil
}
func (t *SegmentSchedule) Join(join event.TournamentEvent) (event.TournamentEvent, error) {
	seg := t.Segments[0]
	join.InstanceId = seg.InstanceId
	total, err := t.updateSegment(join)
	if err != nil {
		return join, err
	}
	seg.TotalJoined = total
	t.TotalJoined += total
	return join, nil
}

func (t *SegmentSchedule) MarshalJSON() ([]byte, error) {
	data := make(map[string]any)
	data["TournamentId"] = fmt.Sprintf("%d", t.TournamentId)
	data["Name"] = t.Name
	data["Schedule"] = t.Schedule
	data["TotalJoined"] = t.TotalJoined
	data["Segments"] = t.Segments
	return json.Marshal(data)
}
