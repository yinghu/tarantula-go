package event

import (
	"fmt"
	"net"
	"strings"

	"gameclustering.com/internal/core"
)

const (
	TCP_READ_BUFFER_SIZE int = 1024
)

type EventEndpoint struct {
	Endpoint string
	Service  EventService
	listener net.Listener

	OutboundEnabled bool

	outboundCQ       chan net.Conn
	outboundEQ       chan Event
	outboundIndex    map[int64]core.DataBuffer
	outboundListener EndpointListener
}

func (s *EventEndpoint) Inbound(client net.Conn, systemId int64) {
	defer func() {
		core.AppLog.Printf("client socket is closed")
		if s.OutboundEnabled {
			ce := CloseEvent{}
			ce.oid = systemId
			s.outboundEQ <- &ce
		}
		client.Close()
	}()
	socket := SocketBuffer{Socket: client, Buffer: make([]byte, TCP_READ_BUFFER_SIZE)}
	for {
		cid, err := socket.ReadInt32()
		if err != nil {
			core.AppLog.Printf("error on read cid %s\n", err.Error())
			s.Service.OnError(nil, err)
			break
		}
		ticket, err := socket.ReadString()
		if err != nil {
			core.AppLog.Printf("error on read ticket %s\n", err.Error())
			s.Service.OnError(nil, err)
			break
		}
		_, err = s.Service.VerifyTicket(ticket)
		if err != nil {
			core.AppLog.Printf("invalid ticket %s\n", ticket)
			s.Service.OnError(nil, err)
			break
		}
		topic, err := socket.ReadString()
		if err != nil {
			core.AppLog.Printf("error on read topic %s\n", err.Error())
			s.Service.OnError(nil, err)
			break
		}
		e, err := s.Service.Create(int(cid), topic)
		if err != nil {
			core.AppLog.Printf("error on create event %s\n", err.Error())
			s.Service.OnError(nil, err)
			break
		}
		err = e.Inbound(&socket)
		if err != nil {
			s.Service.OnError(e, err)
		}
	}
}
func (s *EventEndpoint) Join(client net.Conn) {
	socket := SocketBuffer{Socket: client, Buffer: make([]byte, TCP_READ_BUFFER_SIZE)}
	cid, err := socket.ReadInt32()
	if err != nil {
		core.AppLog.Printf("error on read cid %s\n", err.Error())
		s.Service.OnError(nil, err)
		client.Close()
		return
	}
	if cid != int32(JOIN_CID) {
		core.AppLog.Printf("wrong cid %d\n", cid)
		s.Service.OnError(nil, fmt.Errorf("wrong cid %d", cid))
		client.Close()
		return
	}
	e := JoinEvent{}
	err = e.Inbound(&socket)
	if err != nil {
		s.Service.OnError(&e, err)
		client.Close()
		return
	}
	session, err := s.Service.VerifyTicket(e.Token)
	if err != nil {
		client.Close()
		return
	}
	e.Client = client
	e.Pending = &socket
	e.SystemId = session.SystemId
	s.outboundEQ <- &e
}

func (s *EventEndpoint) Open() error {
	s.outboundIndex = make(map[int64]core.DataBuffer)
	if s.OutboundEnabled {
		s.outboundEQ = make(chan Event, 10)
		s.outboundCQ = make(chan net.Conn, 10)
		go s.outbound()
	}
	parts := strings.Split(s.Endpoint, ":")
	core.AppLog.Printf("Endpoint %s %s\n", parts[0], parts[2])
	server, err := net.Listen(parts[0], ":"+parts[2])
	if err != nil {
		return err
	}
	s.listener = server
	for {
		client, err := s.listener.Accept()
		if err != nil {
			core.AppLog.Printf("Error :%s\n", err.Error())
			break
		}
		if s.OutboundEnabled {
			s.outboundCQ <- client
		} else {
			go s.Inbound(client, 0)
		}
	}
	if s.OutboundEnabled {
		ce := CloseEvent{}
		ce.oid = 0
		s.outboundEQ <- &ce
	}
	core.AppLog.Println("Server closed")
	return nil
}
func (s *EventEndpoint) Close() error {
	core.AppLog.Printf("endpoint shutting down")
	s.listener.Close()
	return nil
}

func (s *EventEndpoint) Push(e Event) {
	s.outboundEQ <- e
}
func (s *EventEndpoint) Register(li EndpointListener) {
	s.outboundListener = li
}
func (s *EventEndpoint) outbound() {
	running := true
	for running {
		select {
		case c := <-s.outboundCQ:
			go s.Join(c)
		case e := <-s.outboundEQ:
			if e.ClassId() == CLOSE_CID {
				if e.OId() == 0 {
					running = false
					continue
				}
				core.AppLog.Printf("remove connection %d\n", e.OId())
				delete(s.outboundIndex, e.OId())
				continue
			}
			if e.ClassId() == JOIN_CID {
				join, _ := e.(*JoinEvent)
				s.outboundIndex[join.SystemId] = join.Pending
				go s.Inbound(join.Client, join.SystemId)
				continue
			}
			s.dispatch(e)
		}
	}
	core.AppLog.Printf("outbound event closed")
}

func (s *EventEndpoint) dispatch(e Event) {
	core.AppLog.Printf("Dispatch event %v\n", e)
	for _, soc := range s.outboundIndex {
		e.Outbound(soc)
	}
}
