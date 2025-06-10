package event

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

const (
	TCP_READ_BUFFER_SIZE int = 1024
)

type Endpoint struct {
	TcpEndpoint string
	Service     EventService
	listener    net.Listener
}

func (s *Endpoint) handleClient(client net.Conn) {
	defer func() {
		client.Close()
	}()
	socket := SocketBuffer{Socket: client, Buffer: make([]byte, TCP_READ_BUFFER_SIZE)}
	cid, err := socket.ReadInt32()
	if err != nil {
		fmt.Printf("Error on read cid %s\n", err.Error())
		return
	}
	e := s.Service.Create(int(cid))
	tik, err := socket.ReadString()
	if err != nil {
		e.OnError(err)
		return
	}
	if tik != "ticket" {
		e.OnError(errors.New("invalid ticket"))
		return
	}
	e.Inbound(&socket)
}

func (s *Endpoint) Open() error {
	parts := strings.Split(s.TcpEndpoint, "://")
	fmt.Printf("Endpoint %s %s\n", parts[0], parts[1])
	server, err := net.Listen(parts[0], parts[1])
	if err != nil {
		return err
	}
	s.listener = server
	for {
		client, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("Error :%s\n", err.Error())
			break
		}
		go s.handleClient(client)
	}
	fmt.Println("Server closed")
	return nil
}
func (s *Endpoint) Close() error {
	s.listener.Close()
	return nil
}
