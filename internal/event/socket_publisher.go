package event

import (
	"net"
	"strings"
)

const SOCKET_READ_BUFFER_SIZE int = 1024


type SocketPublisher struct {
	Remote     string
	BufferSize int
}

func (s *SocketPublisher) Publish(e Event) {
	parts := strings.Split(s.Remote, "://")
	conn, err := net.Dial(parts[0], parts[1])
	if err != nil {
		e.OnError(err)
		return
	}
	defer conn.Close()
	if s.BufferSize==0{
		s.BufferSize = SOCKET_READ_BUFFER_SIZE
	}
	buffer := SocketBuffer{Socket: conn, Buffer: make([]byte, s.BufferSize)}
	err = buffer.WriteInt32(int32(e.ClassId()))
	if err != nil {
		e.OnError(err)
		return
	}
	err = buffer.WriteString("ticket")
	if err != nil {
		e.OnError(err)
		return
	}
	e.Outbound(&buffer)
}
