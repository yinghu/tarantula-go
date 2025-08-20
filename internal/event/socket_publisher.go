package event

import (
	"net"
	"strings"
)

const SOCKET_DATA_BUFFER_SIZE int = 1024

type SocketPublisher struct {
	Remote string
	sb     SocketBuffer
}

func (s *SocketPublisher) Connect() error {
	parts := strings.Split(s.Remote, "://")
	conn, err := net.Dial(parts[0], parts[1])
	if err != nil {
		return err
	}
	s.sb = SocketBuffer{Socket: conn, Buffer: make([]byte, SOCKET_DATA_BUFFER_SIZE)}
	return nil
}

func (s *SocketPublisher) Close() error {
	return s.sb.Clear()
}

func (s *SocketPublisher) Publish(e Event, ticket string) {
	err := s.sb.WriteInt32(int32(e.ClassId()))
	if err != nil {
		e.Listener().OnError(err)
		return
	}
	err = s.sb.WriteString(ticket)
	if err != nil {
		e.Listener().OnError(err)
		return
	}
	err = s.sb.WriteString(e.OnTopic())
	if err != nil {
		e.Listener().OnError(err)
		return
	}
	err = e.Outbound(&s.sb)
	if err != nil {
		e.Listener().OnError(err)
	}
}
