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
		core.AppLog.Printf("outbound message %d\n",e.ClassId())
		out.Soc.WriteInt32(int32(e.ClassId()))
		e.Outbound(out.Soc)
		time.Sleep(10 * time.Millisecond)
	}
}
