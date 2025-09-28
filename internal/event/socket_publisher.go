package event

import (
	"net"
	"strings"
	"time"

	"gameclustering.com/internal/core"
)

const SOCKET_DATA_BUFFER_SIZE int = 1024

type SocketPublisher struct {
	Remote string
	client net.Conn
	pub    SocketBuffer
}

func (s *SocketPublisher) Connect() error {
	parts := strings.Split(s.Remote, "://")
	conn, err := net.Dial(parts[0], parts[1])
	if err != nil {
		return err
	}
	s.client = conn
	s.pub = SocketBuffer{Socket: conn, Buffer: make([]byte, SOCKET_DATA_BUFFER_SIZE)}
	return nil
}

func (s *SocketPublisher) Close() error {
	return s.client.Close()
}

func (s *SocketPublisher) Subscribe(cr EventCreator, ec EventListener) {
	sub := SocketBuffer{Socket: s.client, Buffer: make([]byte, SOCKET_DATA_BUFFER_SIZE)}
	for {
		cid, err := sub.ReadInt32()
		if err != nil {
			ec.OnError(nil, err)
			break
		}
		e, err := cr.Create(int(cid), "")
		if err != nil {
			ec.OnError(nil, err)
			break
		}
		err = e.Inbound(&sub)
		if err != nil {
			ec.OnError(nil, err)
			break
		}
		ec.OnEvent(e)
		time.Sleep(500 * time.Millisecond)
	}
}

func (s *SocketPublisher) Join(e Event) error {
	err := s.pub.WriteInt32(int32(e.ClassId()))
	if err != nil {
		core.AppLog.Printf("error on write classid %s\n", err.Error())
		e.Listener().OnError(e, err)
		return err
	}
	err = e.Outbound(&s.pub)
	if err != nil {
		core.AppLog.Printf("error on write outbound %s\n", err.Error())
		e.Listener().OnError(e, err)
	}
	e.Listener().OnEvent(e)
	return nil
}

func (s *SocketPublisher) Publish(e Event, ticket string) error {
	err := s.pub.WriteInt32(int32(e.ClassId()))
	if err != nil {
		core.AppLog.Printf("error on write classid %s\n", err.Error())
		e.Listener().OnError(e, err)
		return err
	}
	err = s.pub.WriteString(ticket)
	if err != nil {
		core.AppLog.Printf("error on write ticket %s\n", err.Error())
		e.Listener().OnError(e, err)
		return err
	}
	err = s.pub.WriteString(e.Topic())
	if err != nil {
		core.AppLog.Printf("error on write topic %s\n", err.Error())
		e.Listener().OnError(e, err)
		return err
	}
	err = e.Outbound(&s.pub)
	if err != nil {
		core.AppLog.Printf("error on write outbound %s\n", err.Error())
		e.Listener().OnError(e, err)
	}
	e.Listener().OnEvent(e)
	return nil
}
