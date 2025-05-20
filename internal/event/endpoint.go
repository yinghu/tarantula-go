package event

import (
	"fmt"
	"net"
	"strings"
)

type Endpoint struct {
	TcpEndpoint string
	Factory     EventFactory
	listener    net.Listener
}

func (s *Endpoint) handleClient(client net.Conn) {
	defer func() {
		client.Close()
	}()
	socket := SocketBuffer{Socket: client, Buffer: make([]byte, 1024)}
	cid, err := socket.ReadInt()
	if err != nil {
		fmt.Printf("Error on read cid %s\n", err.Error())
		return
	}
	e := s.Factory.Create(cid)
	tik, err := socket.ReadString()
	if err != nil {
		fmt.Printf("Err %s\n", err.Error())
		return
	}
	fmt.Printf("Event : %d %s\n", cid, tik)
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



