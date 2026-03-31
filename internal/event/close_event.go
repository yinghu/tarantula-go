package event

import "gameclustering.com/internal/core"

type CloseEvent struct {
	core.EventObj `json:"-"`
}

func (s *CloseEvent) ClassId() uint32 {
	return CLOSE_CID
}
