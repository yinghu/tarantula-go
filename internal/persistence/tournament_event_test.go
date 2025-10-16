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

type SampleIndexListener struct {
	BadgerLocal
}

func (s *SampleIndexListener) LocalStore() core.DataStore {
	return s
}
func (s *SampleIndexListener) Publish(e event.Event) {
	fmt.Printf("Event : %v\n", e)
}

func TestTournamentEvent(t *testing.T) {
	local := BadgerLocal{InMemory: true, LogDisabled: true}
	err := local.Open()
	if err != nil {
		t.Errorf("Local store error %s", err.Error())
	}
	defer local.Close()
	index := SampleIndexListener{BadgerLocal: local}
	for i := range 5 {
		sid := 1000 + i
		tmnt := event.TournamentEvent{TournamentId: TID, InstanceId: IID, SystemId: int64(sid), Score: 0, LastUpdated: time.Now().UnixMilli()}
		tmnt.LastUpdated = time.Now().UnixMilli()
		err = local.Save(&tmnt)
		if err != nil {
			fmt.Printf("update error %s\n", err.Error())
		}
		tmnt.OnIndex(&index)
	}
	for i := range 5 {
		sid := 1000 + i
		tmnt := event.TournamentEvent{TournamentId: TID, InstanceId: IID, SystemId: int64(sid), Score: 100, LastUpdated: time.Now().UnixMilli()}
		err = local.Load(&tmnt)
		fmt.Printf("Rev : %d\n", tmnt.Revision())
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
	//fmt.Printf("Count : %d %s\n", ct.Count, time.UnixMilli(ct.Timestamp()))
	if ct.Count != 5 {
		t.Errorf("count should be 5 %d", ct.Count)
	}
}

func TestTournamentQuery(t *testing.T) {
	local := BadgerLocal{InMemory: true, LogDisabled: true}
	err := local.Open()
	if err != nil {
		t.Errorf("Local store error %s", err.Error())
	}
	defer local.Close()
	for i := range 10 {
		sid := 10000 + i
		tmnt := event.TournamentEvent{TournamentId: 2000, InstanceId: IID, SystemId: SID, Score: 100, LastUpdated: time.Now().UnixMilli()}
		tmnt.OnOId(int64(sid))
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
	for i := range 10 {
		sid := 20000 + i
		tmnt := event.TournamentEvent{TournamentId: 3000, InstanceId: IID, SystemId: SID, Score: 100, LastUpdated: time.Now().UnixMilli()}
		tmnt.OnOId(int64(sid))
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
	tq := event.QTournament{TournamentId: 2000, InstanceId: 0, SystemId: 0}
	tq.Tag = event.TOURNAMENT_ETAG
	px := BufferProxy{}
	px.NewProxy(100)
	tq.QCriteria(&px)
	px.Flip()
	t2000 := 0
	local.List(&px, func(k, v core.DataBuffer) bool {
		t := event.TournamentEvent{}
		v.ReadInt32()
		v.ReadInt64()
		t.Read(v)
		//t.Rev = rev
		//fmt.Printf("Score %d , LastUpdated %d Rev : %d\n", t.Score, t.LastUpdated, t.Revision())
		t2000++
		return true
	})

	if t2000 != 1 {
		t.Errorf("t2000 should be 1 %d", t2000)
	}

	tq = event.QTournament{TournamentId: 3000, InstanceId: 0, SystemId: 0}
	tq.Tag = event.TOURNAMENT_ETAG
	//px := BufferProxy{}
	//px.NewProxy(100)
	px.Clear()
	tq.QCriteria(&px)
	px.Flip()
	tc := 0
	local.List(&px, func(k, v core.DataBuffer) bool {
		t := event.TournamentEvent{}
		v.ReadInt32()
		v.ReadInt64()
		t.Read(v)
		//t.Rev = rev
		//fmt.Printf("Score %d , LastUpdated %d Rev : %d\n", t.Score, t.LastUpdated, t.Revision())
		tc++
		return true
	})

	if tc != 1 {
		t.Errorf("tc should be 1 %d", tc)
	}

}
