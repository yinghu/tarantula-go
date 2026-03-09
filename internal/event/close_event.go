package event

import "gameclustering.com/internal/core"

type CloseEvent struct {
	core.EventObj `json:"-"`
}

func (s *CloseEvent) ClassId() int {
	return CLOSE_CID
}
