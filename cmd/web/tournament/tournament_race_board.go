package main

import (
	"slices"
	"sync"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type RaceEntry struct {
	SystemId    int64 `json:"SystemId,string"`
	Score       int64 `json:"Score,string"`
	LastUpdated int64 `json:"LastUpdated,string"`
}

type RaceBoard struct {
	sync.RWMutex `json:"-"`
	Size         int32                           `json:"Size"`
	TopListing   []RaceEntry                     `json:"TopListing"`
	index        map[int64]event.TournamentEvent `json:"-"`
}

func cmp(a, b RaceEntry) int {
	diff := a.Score - b.Score
	if diff > 0 {
		return -1
	}
	if diff < 0 {
		return 1
	}
	tm := a.LastUpdated - b.LastUpdated
	if tm <= 0 {
		return -1
	}
	return 1
}

func (b *RaceBoard) Start() {
	b.index = make(map[int64]event.TournamentEvent)
	b.TopListing = make([]RaceEntry, b.Size)
	for i := range b.TopListing {
		b.TopListing[i] = RaceEntry{SystemId: int64(i)}
	}
	slices.SortFunc(b.TopListing, cmp)
}

func (b *RaceBoard) OnBoard(te event.TournamentEvent) {
	core.AppLog.Printf("on board %v\n", te)
	b.Lock()
	defer b.Unlock()
	if te.Score < b.TopListing[b.Size-1].Score {
		return
	}
	delete(b.index, b.TopListing[b.Size-1].SystemId)
	b.index[te.SystemId] = te
	ix := 0
	for k, v := range b.index {
		b.TopListing[ix] = RaceEntry{SystemId: k, Score: v.Score, LastUpdated: v.LastUpdated}
		ix++
	}
	slices.SortFunc(b.TopListing, cmp)
}

func (b *RaceBoard) Listing() []RaceEntry {
	core.AppLog.Printf("on listing")
	b.RLock()
	defer b.RUnlock()
	listing := make([]RaceEntry, 0)
	for i := range b.TopListing {
		re := b.TopListing[i]
		if re.SystemId < 100 {
			continue
		}
		listing = append(listing, re)
	}
	return listing
}
