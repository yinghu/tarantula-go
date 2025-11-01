package event

import (
	"net"
	"strings"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/metrics"
)

type TcpEndpoint struct {
	Endpoint string
	Service  EventService
	listener net.Listener

	OutboundEnabled bool

	outboundCQ    chan net.Conn
	outboundEQ    chan Event
	outboundIndex map[int64]*OutboundSocket
}

func (s *TcpEndpoint) Open() error {
	s.outboundIndex = make(map[int64]*OutboundSocket)
	if s.OutboundEnabled {
		s.outboundEQ = make(chan Event, 10)
		s.outboundCQ = make(chan net.Conn, 10)
		go s.outbound()
	}
	parts := strings.Split(s.Endpoint, ":")
	core.AppLog.Printf("Endpoint %s :%s\n", parts[0], parts[2])
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
			go s.inbound(client, 0)
		}
	}
	if s.OutboundEnabled {
		ce := CloseEvent{}
		ce.EOid = 0
		s.outboundEQ <- &ce
	}
	core.AppLog.Println("Server closed")
	close(s.outboundCQ)
	close(s.outboundEQ)
	return nil
}

func (s *TcpEndpoint) inbound(client net.Conn, systemId int64) {

}

func (s *TcpEndpoint) outbound() {
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
					metrics.SOCKET_CONCURRENCY_METRICS.Dec()
				}
				continue
			}
			if e.ClassId() == JOIN_CID {
				metrics.SOCKET_CONCURRENCY_METRICS.Inc()
				join, _ := e.(*JoinEvent)
				cout := OutboundSocket{Soc: join.Pending, Pending: make(chan Event, 10)}
				go cout.Subscribe()
				s.outboundIndex[join.SystemId] = &cout
				go s.inbound(join.Client, join.SystemId)
				continue
			}
			s.dispatch(e)
		}
	}
	core.AppLog.Printf("outbound event closed")
}

func (s *TcpEndpoint) join(client net.Conn) {

}

func (s *TcpEndpoint) dispatch(e Event) {

}
