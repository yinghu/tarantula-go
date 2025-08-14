package persistence

import (
	"fmt"
	"testing"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

const (
	TID int64 = 100
	IID int64 = 200
	SID int64 = 300
)

func TestTournamentEvent(t *testing.T) {
	local := BadgerLocal{InMemory: false, Path: "/home/yinghu/local/test"}
	err := local.Open()
	if err != nil {
		t.Errorf("Local store error %s", err.Error())
	}
	defer local.Close()
	for i := range 5 {
		sid := 1000 + i
		tmnt := event.TournamentEvent{Id: int64(sid), TournamentId: TID, InstanceId: IID, SystemId: SID, Score: 100, LastUpdated: time.Now().UnixMilli()}
		err = local.Load(&tmnt)
		if err != nil { //not fount
			err = local.Save(&tmnt)
			if err != nil {
				fmt.Printf("new save error %s\n", err.Error())
			}
		} else {
			tmnt.Score = tmnt.Score + 100
			tmnt.LastUpdated = time.Now().UnixMilli()
			err = local.Save(&tmnt)
			if err != nil {
				fmt.Printf("update error %s\n", err.Error())
			}
		}
	}
	for i := range 5 {
		sid := 1000 + i
		tmnt := event.TournamentEvent{Id: int64(sid), TournamentId: TID, InstanceId: IID, SystemId: SID, Score: 100, LastUpdated: time.Now().UnixMilli()}
		err = local.Load(&tmnt)
		if err != nil { //not fount
			err = local.Save(&tmnt)
			if err != nil {
				fmt.Printf("new save error %s\n", err.Error())
			}
		} else {
			tmnt.Score = tmnt.Score + 100
			tmnt.LastUpdated = time.Now().UnixMilli()
			err = local.Save(&tmnt)
			if err != nil {
				fmt.Printf("update error %s\n", err.Error())
			}
		}
	}
	ct := event.StatEvent{Tag: event.TOURNAMENT_ETAG, Name: event.STAT_TOTAL}
	err = local.Load(&ct)
	if err != nil {
		t.Errorf("Local store error %s", err.Error())
	}
	fmt.Printf("Count : %d\n", ct.Count)
	if ct.Count != 5 {
		t.Errorf("count should be 5 %d", ct.Count)
	}
	px := BufferProxy{}
	px.NewProxy(100)
	px.WriteString(event.TOURNAMENT_ETAG)
	px.WriteInt64(TID)
	px.WriteInt64(IID)
	px.WriteInt64(SID)
	px.Flip()
	t1000 := 0
	local.List(&px, func(k, v core.DataBuffer, rev uint64) bool {
		t := event.TournamentEvent{}
		v.ReadInt32()
		t.Read(v)
		t.Rev = rev
		fmt.Printf("Score %d , LastUpdated %d Rev : %d\n", t.Score, t.LastUpdated, t.Revision())
		t1000++
		return true
	})

	if t1000 != 5 {
		t.Errorf("t200 should be 5 %d", t1000)
	}

}
