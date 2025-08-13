package main

import (
	"encoding/json"
	"fmt"

	"gameclustering.com/internal/event"
)

type Segment struct {
	InstanceId int64  `json:"-"`
	Name       string `json:"Name"`
}

type SegementSchedule struct {
	TournamentId int64  `json:"TournamentId,string"`
	Name         string `json:"Name"`
	TotalJoined  int32 `json:"TotalJoined"`
	Schedule
	Segments []Segment `json:"-"`

	*TournamentService `json:"-"`
}

func (t *SegementSchedule) Join(join event.TournamentEvent) error {
	seg := t.Segments[0]
	join.InstanceId = seg.InstanceId
	total, err := t.updateSegment(join)
	if err != nil {
		return err
	}
	t.TotalJoined = total
	return nil
}

func (t *SegementSchedule) MarshalJSON() ([]byte, error) {
	data := make(map[string]any)
	data["TournamentId"] = fmt.Sprintf("%d", t.TournamentId)
	data["Name"] = t.Name
	data["Schedule"] = t.Schedule
	data["TotalJoined"] = t.TotalJoined
	return json.Marshal(data)
}
