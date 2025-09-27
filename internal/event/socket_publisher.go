package event

import (
	"net"
	"strings"

	"gameclustering.com/internal/core"
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

func (s *SocketPublisher) Join(e Event) error {
	err := s.sb.WriteInt32(int32(e.ClassId()))
	if err != nil {
		core.AppLog.Printf("error on write classid %s\n", err.Error())
		e.Listener().OnError(e, err)
		return err
	}
	err = e.Outbound(&s.sb)
	if err != nil {
		core.AppLog.Printf("error on write outbound %s\n", err.Error())
		e.Listener().OnError(e, err)
	}
	e.Listener().OnEvent(e)
	return nil
}

func (s *SocketPublisher) Publish(e Event, ticket string) error {
	err := s.sb.WriteInt32(int32(e.ClassId()))
	if err != nil {
		core.AppLog.Printf("error on write classid %s\n", err.Error())
		e.Listener().OnError(e, err)
		return err
	}
	err = s.sb.WriteString(ticket)
	if err != nil {
		core.AppLog.Printf("error on write ticket %s\n", err.Error())
		e.Listener().OnError(e, err)
		return err
	}
	err = s.sb.WriteString(e.Topic())
	if err != nil {
		core.AppLog.Printf("error on write topic %s\n", err.Error())
		e.Listener().OnError(e, err)
		return err
	}
	err = e.Outbound(&s.sb)
	if err != nil {
		core.AppLog.Printf("error on write outbound %s\n", err.Error())
		e.Listener().OnError(e, err)
	}
	e.Listener().OnEvent(e)
	return nil
}
