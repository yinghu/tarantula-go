package event

import (
	"net"
	"strings"
)

type SocketPublisher struct {
	Remote     string
	BufferSize int
}

func (s *SocketPublisher) Publish(e Event) error {
	parts := strings.Split(s.Remote, "://")
	conn, err := net.Dial(parts[0], parts[1])
	if err != nil {
		return err
	}
	defer conn.Close()
	buffer := SocketBuffer{Socket: conn, Buffer: make([]byte, s.BufferSize)}
	buffer.WriteInt32(int32(e.ClassId()))
	buffer.WriteString("ticket")
	e.Outbound(&buffer)
	return nil
}
