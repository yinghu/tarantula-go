package event

import (
	"net"
	"strings"
	"time"

	"gameclustering.com/internal/core"
)

//const SOCKET_DATA_BUFFER_SIZE int = 1024

type TcpPublisher struct {
	Remote string
	client net.Conn
	pub    core.DataBuffer
}

func (s *TcpPublisher) Connect() error {
	parts := strings.Split(s.Remote, "://")
	conn, err := net.Dial(parts[0], parts[1])
	if err != nil {
		return err
	}
	s.client = conn
	s.pub = core.NewBuffer(SOCKET_DATA_BUFFER_SIZE)
	return nil
}

func (s *TcpPublisher) Close() error {
	return s.client.Close()
}

func (s *TcpPublisher) Subscribe(cr EventCreator, ec EventListener) {
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
		time.Sleep(10 * time.Millisecond)
	}
}

func (s *TcpPublisher) Join(e Event) error {
	s.pub.Clear()
	err := s.pub.WriteInt32(int32(e.ClassId()))
	if err != nil {
		core.AppLog.Printf("error on write classid %s\n", err.Error())
		e.Listener().OnError(e, err)
		return err
	}
	err = e.Outbound(s.pub)
	if err != nil {
		core.AppLog.Printf("error on write outbound %s\n", err.Error())
		e.Listener().OnError(e, err)
	}
	s.pub.Flip()
	data, err := s.pub.Export('|')
	if err != nil {
		core.AppLog.Printf("error on export %s\n", err.Error())
		e.Listener().OnError(e, err)
	}
	n, err := s.client.Write(data)
	if err != nil {
		core.AppLog.Printf("error on write socket %s\n", err.Error())
		e.Listener().OnError(e, err)
	}
	core.AppLog.Printf("write socket number %d\n", n)
	//e.Listener().OnEvent(e)
	return nil
}

func (s *TcpPublisher) Publish(e Event, ticket string) error {
	s.pub.Clear()
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
	err = e.Outbound(s.pub)
	if err != nil {
		core.AppLog.Printf("error on write outbound %s\n", err.Error())
		e.Listener().OnError(e, err)
	}
	s.pub.Flip()
	data, err := s.pub.Export('|')
	if err != nil {
		core.AppLog.Printf("error on export %s\n", err.Error())
		e.Listener().OnError(e, err)
	}
	n, err := s.client.Write(data)
	if err != nil {
		core.AppLog.Printf("error on write socket %s\n", err.Error())
		e.Listener().OnError(e, err)
	}
	core.AppLog.Printf("write socket number %d\n", n)
	//e.Listener().OnEvent(e)
	return nil
}
