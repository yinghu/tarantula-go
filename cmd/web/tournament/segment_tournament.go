package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"gameclustering.com/internal/event"
)

type Segment struct {
	InstanceId  int64  `json:"-"`
	Name        string `json:"Name"`
	TotalJoined int32  `json:"TotalJoined"`
	RaceBoard   `json:"-"`
}

type SegmentSchedule struct {
	TournamentId int64  `json:"TournamentId,string"`
	Name         string `json:"Name"`
	TotalJoined  int32  `json:"TotalJoined"`
	Schedule
	Segments []*Segment `json:"Segments"`

	*TournamentService `json:"-"`

	sync.RWMutex `json:"-"`
	Started      bool `json:"-"`
}

func (t *SegmentSchedule) Start() error {
	t.Lock()
	defer t.Unlock()
	for i := range t.Segments {
		if !t.Started {
			t.Segments[i].RaceBoard = RaceBoard{Size: 16}
			t.Segments[i].RaceBoard.Start()
		}
		tq := event.QTournament{TournamentId: t.TournamentId, InstanceId: t.Segments[i].InstanceId}
		tq.Tag = event.TOURNAMENT_ETAG
		tq.Id = event.Q_TOURNAMENT_QID
		tq.Topic = "tournament"
		t.Recover(&tq)
	}
	t.Started = true
	return nil
}

func (t *SegmentSchedule) Score(score event.TournamentEvent) (event.TournamentEvent, error) {
	score.LastUpdated = time.Now().UnixMilli()
	sc, err := t.updateEntry(score)
	if err != nil {
		return score, err
	}
	score.Score = sc
	go t.sendEvent(score)
	return score, nil
}
func (t *SegmentSchedule) Join(join event.TournamentEvent) (event.TournamentEvent, error) {
	joined := t.checkJoin(join)
	if joined.InstanceId > 0 {
		return joined, nil
	}
	seg := t.Segments[0]
	join.InstanceId = seg.InstanceId
	total, err := t.updateSegment(join)
	if err != nil {
		return join, err
	}
	seg.TotalJoined = total
	t.TotalJoined += total
	join.LastUpdated = time.Now().UnixMilli()
	go t.sendEvent(join)
	return join, nil
}
func (t *SegmentSchedule) OnBoard(update event.TournamentEvent) {
	t.Segments[0].OnBoard(update)
}
func (t *SegmentSchedule) Listing(te event.TournamentEvent) []RaceEntry {
	return t.Segments[0].Listing()
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

func (t *SegmentSchedule) sendEvent(te event.TournamentEvent) {
	id, _ := t.Sequence().Id()
	e := event.TournamentEvent{TournamentId: te.TournamentId, InstanceId: te.InstanceId, SystemId: te.SystemId, Score: te.Score, LastUpdated: te.LastUpdated}
	e.OId(id)
	e.Topic("tournament")
	t.Send(&e)
}
