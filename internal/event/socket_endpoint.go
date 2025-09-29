package event

import (
	"fmt"
	"net"
	"strings"
	"time"

	"gameclustering.com/internal/core"
)

const (
	TCP_READ_BUFFER_SIZE int    = 1024
	JOIN_TOPIC           string = "join"
)

type SocketEndpoint struct {
	Endpoint string
	Service  EventService
	listener net.Listener

	OutboundEnabled bool

	outboundCQ    chan net.Conn
	outboundEQ    chan Event
	outboundIndex map[int64]*OutboundSocket
}

func (s *SocketEndpoint) Inbound(client net.Conn, systemId int64) {
	defer func() {
		core.AppLog.Printf("client socket is closed")
		if s.OutboundEnabled {
			ce := KickoffEvent{}
			ce.SystemId = systemId
			ce.Source = "disconnected"
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
		time.Sleep(10 * time.Millisecond)
	}
}
func (s *SocketEndpoint) join(client net.Conn) {
	time.Sleep(10 * time.Millisecond)
	socket := SocketBuffer{Socket: client, Buffer: make([]byte, TCP_READ_BUFFER_SIZE)}
	cid, err := socket.ReadInt32()
	if err != nil {
		core.AppLog.Printf("error on read cid %s\n", err.Error())
		s.Service.OnError(nil, err)
		client.Close()
		return
	}
	e, err := s.Service.Create(int(cid), JOIN_TOPIC)
	if err != nil {
		core.AppLog.Printf("wrong join cid %d\n", cid)
		s.Service.OnError(nil, fmt.Errorf("wrong join cid %d", cid))
		client.Close()
		return
	}
	err = e.Inbound(&socket)
	if err != nil {
		s.Service.OnError(e, err)
		client.Close()
		return
	}
	join, ok := e.(*JoinEvent)
	if !ok {
		core.AppLog.Printf("wrong join event %d\n", cid)
		s.Service.OnError(nil, fmt.Errorf("wrong join cid %d", cid))
		client.Close()
		return
	}
	session, err := s.Service.VerifyTicket(join.Ticket)
	if err != nil {
		core.AppLog.Printf("wrong permission %s\n", err.Error())
		client.Close()
		return
	}
	join.Client = client
	join.Pending = &socket
	join.SystemId = session.SystemId
	s.Service.OnEvent(join)
	s.outboundEQ <- join
}

func (s *SocketEndpoint) Open() error {
	s.outboundIndex = make(map[int64]*OutboundSocket)
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
	close(s.outboundCQ)
	close(s.outboundEQ)
	return nil
}
func (s *SocketEndpoint) Close() error {
	core.AppLog.Printf("endpoint shutting down")
	s.listener.Close()
	return nil
}

func (s *SocketEndpoint) Push(e Event) {
	s.outboundEQ <- e
}

func (s *SocketEndpoint) outbound() {
	running := true
	for running {
		select {
		case c := <-s.outboundCQ:
			go s.join(c)
		case e := <-s.outboundEQ:
			if e.ClassId() == CLOSE_CID {
				running = false
				continue
			}
			if e.ClassId() == KICKOFF_CID {
				oc, exists := s.outboundIndex[e.RecipientId()]
				if exists {
					core.AppLog.Printf("remove connection from %d\n", e.RecipientId())
					close(oc.Pending)
					delete(s.outboundIndex, e.RecipientId())
					s.Service.OnEvent(e)
				}
				continue
			}
			if e.ClassId() == JOIN_CID {
				join, _ := e.(*JoinEvent)
				cout := OutboundSocket{Soc: join.Pending, Pending: make(chan Event, 10)}
				go cout.Subscribe()
				s.outboundIndex[join.SystemId] = &cout
				go s.Inbound(join.Client, join.SystemId)
				continue
			}
			s.dispatch(e)
		}
	}
	core.AppLog.Printf("outbound event closed")
}

func (s *SocketEndpoint) dispatch(e Event) {
	core.AppLog.Printf("dispatching %v\n", e)
	if e.RecipientId() > 0 {
		soc, exists := s.outboundIndex[e.RecipientId()]
		core.AppLog.Printf("dispatching %v %v\n", e, exists)
		if exists {
			soc.Pending <- e
		}
		return
	}
	core.AppLog.Printf("all dispatching %v\n", e)
	for _, soc := range s.outboundIndex {
		soc.Pending <- e
	}
}
