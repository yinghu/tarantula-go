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
	outboundIndex map[int64]*OutboundSoc
}

func (s *TcpEndpoint) Open() error {
	s.outboundIndex = make(map[int64]*OutboundSoc)
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
	data := make([]byte, TCP_READ_BUFFER_SIZE)
	buff := core.NewBuffer(TCP_READ_BUFFER_SIZE)
	for {
		num, err := client.Read(data)
		if err != nil {
			core.AppLog.Printf("close inbound")
			s.Service.OnError(nil, err)
			client.Close()
			return
		}
		core.AppLog.Printf("RC %d\n", num)
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
		tick, err := buff.ReadString()
		if err != nil {
			buff.Clear()
			continue
		}
		topic, err := buff.ReadString()
		if err != nil {
			buff.Clear()
			continue
		}
		core.AppLog.Printf("%d %s %s\n", cid, tick, topic)
		e, err := s.Service.Create(int(cid), topic)
		if err != nil {
			buff.Clear()
			continue
		}
		e.Inbound(buff)
		buff.Clear()
		e.Listener().OnEvent(e)
	}
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
				core.AppLog.Printf("KICK OFF")
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
				cout := OutboundSoc{C: join.Client, Pending: make(chan Event, 10)}
				go cout.Sub()
				s.outboundIndex[join.SystemId] = &cout
				go s.inbound(join.Client, join.SystemId)
				continue
			}
			s.dispatch(e)
		}
	}
	core.AppLog.Printf("outbound event closed")
}
func (s *TcpEndpoint) Push(e Event) {
	s.outboundEQ <- e
}
func (s *TcpEndpoint) join(client net.Conn) {
	data := make([]byte, TCP_READ_BUFFER_SIZE)
	buff := core.NewBuffer(TCP_READ_BUFFER_SIZE)
	for {
		num, err := client.Read(data)
		if err != nil {
			s.Service.OnError(nil, err)
			client.Close()
			return
		}
		core.AppLog.Printf("RC %d\n", num)
		err = buff.Write(data[:num])
		if err != nil {
			core.AppLog.Printf("write buff error %s\n", err.Error())
			return
		}
		if data[num-1] == '|' {
			break
		}
	}
	buff.Flip()
	buff.ReadInt32()
	e := JoinEvent{}
	e.Inbound(buff)
	session, err := s.Service.VerifyTicket(e.Ticket)
	if err != nil {
		core.AppLog.Printf("wrong permission %s\n", err.Error())
		client.Close()
		return
	}
	e.Client = client
	e.SystemId = session.SystemId
	s.outboundEQ <- &e
}

func (s *TcpEndpoint) dispatch(e Event) {
	if e.RecipientId() > 0 {
		soc, exists := s.outboundIndex[e.RecipientId()]
		if exists {
			soc.Pending <- e
		}
		return
	}
	for _, soc := range s.outboundIndex {
		soc.Pending <- e
	}
}
