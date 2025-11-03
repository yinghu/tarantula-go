package event

import (
	"net"
	"time"

	"gameclustering.com/internal/core"
)



type OutboundSocket struct {
	C       net.Conn
	B       core.DataBuffer
	Pending chan Event
}

func (out *OutboundSocket) Subscribe() {
	for e := range out.Pending {
		out.B.Clear()
		err := out.B.WriteInt32(int32(e.ClassId()))
		if err != nil {
			break
		}
		err = e.Outbound(out.B)
		if err != nil {
			break
		}
		out.B.Flip()
		data, err := out.B.Export('|')
		if err != nil {
			break
		}
		num, err := out.C.Write(data)
		if err != nil {
			break
		}
		core.AppLog.Printf("write number %d\n", num)
		time.Sleep(10 * time.Millisecond)
	}
}
