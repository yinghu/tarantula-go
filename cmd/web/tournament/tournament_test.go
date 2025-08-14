package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"gameclustering.com/internal/event"
)

func TestToJson(t *testing.T) {
	mp := make(map[int64]Tournament)
	seg := SegmentSchedule{TournamentId: 100}
	mp[seg.TournamentId] = &seg
	d, err := json.Marshal(mp)
	if err != nil {
		t.Errorf("Service error %s", err.Error())
	}
	fmt.Printf("D : %s\n", string(d))
}

func TestToRaceBoard(t *testing.T) {
	rb := RaceBoard{Size: 10}
	rb.Start()
	l0 := rb.Listing()
	z0 := len(l0)
	if z0 != 0 {
		t.Errorf("should be zero length %d", z0)
	}
	tm := time.Now().UnixMilli()
	te1 := event.TournamentEvent{SystemId: 300, Score: 100, LastUpdated: tm}
	rb.OnBoard(te1)
	l1 := rb.Listing()
	z1 := len(l1)
	if z1 != 1 {
		t.Errorf("should be one length %d", z1)
	}

	te2 := event.TournamentEvent{SystemId: 301, Score: 100, LastUpdated: tm + 1}
	rb.OnBoard(te2)
	l2 := rb.Listing()
	z2 := len(l2)
	if z2 != 2 {
		t.Errorf("should be two length %d", z2)
	}
	if l2[0].SystemId != 300 {
		t.Errorf("first should be 300 %d", l2[0].SystemId)
	}

	te3 := event.TournamentEvent{SystemId: 302, Score: 100, LastUpdated: tm + 2}
	rb.OnBoard(te3)
	l3 := rb.Listing()
	z3 := len(l3)
	if z3 != 3 {
		t.Errorf("should be three length %d", z3)
	}
	if l3[2].SystemId != 302 {
		t.Errorf("third should be 302 %d", l3[2].SystemId)
	}
	var ts int64 = 10
	for i := range 100 {
		sid := 1000 + i
		ts += tm
		te := event.TournamentEvent{SystemId: int64(sid), Score: int64(sid), LastUpdated: ts}
		rb.OnBoard(te)
	}
	l10 := rb.Listing()
	z10 := len(l10)
	if z10 != 10 {
		t.Errorf("should be ten length %d", z10)
	}

}
