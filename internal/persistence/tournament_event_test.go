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
	endTime := time.Now().Add(1 * time.Hour).UnixMilli()
	for i := range 10 {
		sid := 1000 + i
		tmnt := event.TournamentEvent{TournamentId: TID, InstanceId: IID, SystemId: int64(sid), Score: 0, LastUpdated: endTime - time.Now().UnixMilli()}
		tmnt.LastUpdated = endTime - time.Now().UnixMilli()
		err = local.Save(&tmnt)
		if err != nil {
			fmt.Printf("update error %s\n", err.Error())
		}
		tmnt.OnIndex(&index)
	}
	time.Sleep(100 * time.Millisecond)
	for i := range 10 {
		sid := 1000 + i
		tmnt := event.TournamentEvent{TournamentId: TID, InstanceId: IID, SystemId: int64(sid), Score: 100, LastUpdated: endTime - time.Now().UnixMilli()}
		err = local.Load(&tmnt)
		if err != nil { //not fount
			err = local.Save(&tmnt)
			if err != nil {
				fmt.Printf("new save error %s\n", err.Error())
			}
		} else {
			tmnt.Score = tmnt.Score + 100 + int64(i)
			tmnt.LastUpdated = endTime - time.Now().UnixMilli()
			err = local.Save(&tmnt)
			if err != nil {
				fmt.Printf("update error %s\n", err.Error())
			}
			tmnt.OnIndex(&index)
		}
	}
	time.Sleep(100 * time.Millisecond)
	for i := range 10 {
		sid := 1000 + i
		tmnt := event.TournamentEvent{TournamentId: TID, InstanceId: IID, SystemId: int64(sid), Score: 100, LastUpdated: endTime - time.Now().UnixMilli()}
		err = local.Load(&tmnt)
		if err != nil { //not fount
			err = local.Save(&tmnt)
			if err != nil {
				fmt.Printf("new save error %s\n", err.Error())
			}
		} else {
			tmnt.Score = tmnt.Score + 200 + int64(i)
			tmnt.LastUpdated = endTime - time.Now().UnixMilli()
			err = local.Save(&tmnt)
			if err != nil {
				fmt.Printf("update error %s\n", err.Error())
			}
			tmnt.OnIndex(&index)
		}
	}
	tq := event.QScore{TournamentId: TID, InstanceId: IID}
	prefix := NewBuffer(100)
	tq.QCriteria(prefix)
	prefix.Flip()
	kp, _ := prefix.Read(0)
	opt := core.ListingOpt{Prefix: kp, Reverse: true, Limit: 0}
	imp := make(map[int64]event.TournamentScoreIndex)
	err = local.Query(opt, func(k, v core.DataBuffer) bool {
		tc := event.TournamentScoreIndex{}
		tc.ReadKey(k)
		_, exists := imp[tc.SystemId]
		if exists {
			return true
		}
		imp[tc.SystemId] = tc
		fmt.Printf("TID : %d , INS : %d , SCORE : %d , TM : %d , SYS : %d\n", tc.TournamentId, tc.InstanceId, tc.Score, tc.UpdateTime, tc.SystemId)
		return true
	})
	if err != nil {
		t.Errorf("should be error %s", err.Error())
	}
}
