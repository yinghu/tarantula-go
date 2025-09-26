package event

import (
	"net"
	"strings"
	"sync"

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

	lock          sync.Mutex
	outboundQueue chan Event
	outboundIndex map[string]core.DataBuffer
}

func (s *EventEndpoint) Inbound(client net.Conn) {
	defer func() {
		core.AppLog.Printf("client socket is closed")
		if s.OutboundEnabled {
			s.removeOutbound(client)
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
		err = s.Service.VerifyTicket(ticket)
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

func (s *EventEndpoint) addOutbound(client net.Conn) {
	s.lock.Lock()
	defer s.lock.Unlock()
	soc := SocketBuffer{Socket: client, Buffer: make([]byte, TCP_READ_BUFFER_SIZE)}
	cid := client.RemoteAddr().String()
	s.outboundIndex[cid] = &soc
	core.AppLog.Printf("client added %s\n", cid)
}
func (s *EventEndpoint) removeOutbound(client net.Conn) {
	s.lock.Lock()
	defer s.lock.Unlock()
	cid := client.RemoteAddr().String()
	delete(s.outboundIndex, cid)
	core.AppLog.Printf("client removed %s\n", cid)
}

func (s *EventEndpoint) Open() error {
	s.lock = sync.Mutex{}
	s.outboundIndex = make(map[string]core.DataBuffer)
	if s.OutboundEnabled {
		s.outboundQueue = make(chan Event, 10)
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
		go s.Inbound(client)
		if s.OutboundEnabled {
			s.addOutbound(client)
		}
	}
	if s.OutboundEnabled {
		s.outboundQueue <- &CloseEvent{}
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
	s.outboundQueue <- e
}
func (s *EventEndpoint) outbound() {
	for e := range s.outboundQueue {
		if e.ClassId() == CLOSE_CID {
			break
		}
		s.dispatch(e)
	}
	core.AppLog.Printf("outbound event closed")
}

func (s *EventEndpoint) dispatch(e Event) {
	core.AppLog.Printf("Dispatch event %v\n", e)
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, soc := range s.outboundIndex {
		e.Outbound(soc)
	}
}
