package event

import (
	"fmt"
	"net"
	"strings"
)

type Endpoint struct {
	tcpEndpoint string
	listener    net.Listener
}

func (s *Endpoint) handleClient(client net.Conn) {
	defer func() {
		client.Close()
		s.Close()
	}()
	reader := SocketReader{Socket: client, Buffer: make([]byte, 1024)}
	fid, _ := reader.ReadInt32()
	cid, _ := reader.ReadInt32()
	tik, err := reader.ReadString()
	if err != nil {
		fmt.Printf("Err %s\n", err.Error())
	}
	fmt.Printf("HS %d : %d : %s\n", fid, cid, tik)
}

func (s *Endpoint) Open() error {
	parts := strings.Split(s.tcpEndpoint, "://")
	fmt.Printf("Endpoint %s %s\n", parts[0], parts[1])
	server, err := net.Listen(parts[0], parts[1])
	if err != nil {
		return err
	}
	s.listener = server
	for {
		client, er := s.listener.Accept()
		if er != nil {
			fmt.Printf("Error :%s\n", er.Error())
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

func (s *Endpoint) Publish(event Event) error {
	event.Streaming(Chunk{true, []byte("hellop")})
	event.Streaming(Chunk{true, []byte("hellop")})
	event.Streaming(Chunk{false, []byte("hellop")})
	return nil
}
