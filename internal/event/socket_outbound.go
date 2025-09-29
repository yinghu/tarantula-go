package event

import (
	"time"

	"gameclustering.com/internal/core"
)

type OutboundSocket struct {
	Soc     core.DataBuffer
	Pending chan Event
}

func (out *OutboundSocket) Subscribe() {
	for e := range out.Pending {
		err := out.Soc.WriteInt32(int32(e.ClassId()))
		if err != nil {
			break
		}
		err = e.Outbound(out.Soc)
		if err != nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}
