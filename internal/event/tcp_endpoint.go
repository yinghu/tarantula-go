package event

import (
	"net"
	"strings"

	"gameclustering.com/internal/core"
)

const (
	TCP_READ_BUFFER_SIZE int = 1024
)

type TcpEndpoint struct {
	Endpoint string
	Service  EventService
	listener net.Listener
}

func (s *TcpEndpoint) handleClient(client net.Conn) {
	defer func() {
		client.Close()
	}()
	socket := SocketBuffer{Socket: client, Buffer: make([]byte, TCP_READ_BUFFER_SIZE)}
	running := true
	for running {
		cid, err := socket.ReadInt32()
		if err != nil {
			core.AppLog.Printf("error on read cid %s\n", err.Error())
			s.Service.OnError(err)
			break
		}
		ticket, err := socket.ReadString()
		if err != nil {
			core.AppLog.Printf("error on read ticket %s\n", err.Error())
			s.Service.OnError(err)
			break
		}
		err = s.Service.VerifyTicket(ticket)
		if err != nil {
			core.AppLog.Printf("invalid ticket %s\n", ticket)
			s.Service.OnError(err)
			break
		}
		topic, err := socket.ReadString()
		if err != nil {
			core.AppLog.Printf("error on read topic %s\n", err.Error())
			s.Service.OnError(err)
			break
		}
		e, err := s.Service.Create(int(cid), topic)
		if err != nil {
			core.AppLog.Printf("error on create event %s\n", err.Error())
			s.Service.OnError(err)
			break
		}
		err = e.Inbound(&socket)
		if err != nil {
			s.Service.OnError(err)
		}
	}
}

func (s *TcpEndpoint) Open() error {
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
		go s.handleClient(client)
	}
	core.AppLog.Println("Server closed")
	return nil
}
func (s *TcpEndpoint) Close() error {
	s.listener.Close()
	return nil
}
