package event

import (
	"net"
	"strings"

	"gameclustering.com/internal/core"
)


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
	s.pub = core.NewBuffer(TCP_READ_BUFFER_SIZE)
	return nil
}

func (s *TcpPublisher) Close() error {
	return s.client.Close()
}

func (s *TcpPublisher) Subscribe(cr EventCreator, ec EventListener) {
	buff := core.NewBuffer(TCP_READ_BUFFER_SIZE)
	data := make([]byte, TCP_READ_BUFFER_SIZE)
	for {
		num, err := s.client.Read(data)
		if err != nil {
			ec.OnError(nil, err)
			s.client.Close()
			return
		}
		err = buff.Write(data[:num])
		if err != nil {
			core.AppLog.Printf("write buff error %s\n", err.Error())
			return
		}
		if data[num-1] != '|' {
			continue
		}
		buff.Flip()
		cid, err := buff.ReadInt32()
		if err != nil {
			buff.Clear()
			continue
		}
		e, err := cr.Create(int(cid),"local")
		if err != nil {
			buff.Clear()
			continue
		}
		e.Inbound(buff)
		buff.Clear()
		e.Listener().OnEvent(e)
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
	_, err = s.client.Write(data)
	if err != nil {
		core.AppLog.Printf("error on write socket %s\n", err.Error())
		e.Listener().OnError(e, err)
	}
	
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
	_, err = s.client.Write(data)
	if err != nil {
		core.AppLog.Printf("error on write socket %s\n", err.Error())
		e.Listener().OnError(e, err)
	}
	return nil
}
