package main

import (
	"encoding/json"
	"fmt"
	"testing"
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
