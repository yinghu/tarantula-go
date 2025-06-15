package event

import (
	"fmt"
	"net"
	"strings"
)

const SOCKET_DATA_BUFFER_SIZE int = 1024

type SocketPublisher struct {
	Remote string
}

func (s *SocketPublisher) Publish(e Event) {
	fmt.Printf("publish %s\n", s.Remote)
	parts := strings.Split(s.Remote, "://")
	conn, err := net.Dial(parts[0], parts[1])
	if err != nil {
		e.OnError(err)
		return
	}
	defer conn.Close()

	buffer := SocketBuffer{Socket: conn, Buffer: make([]byte, SOCKET_DATA_BUFFER_SIZE)}
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
