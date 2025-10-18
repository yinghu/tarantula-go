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
func (s *SampleIndexListener) Index(e event.Event) {
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
		tmnt := event.TournamentEvent{TournamentId: TID, InstanceId: IID, SystemId: int64(sid), Score:0, LastUpdated: endTime - time.Now().UnixMilli()}
		tmnt.LastUpdated = endTime - time.Now().UnixMilli()
		tmnt.OnIndex(&index)
	}
}
