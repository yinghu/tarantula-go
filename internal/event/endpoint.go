package event

import (
	"fmt"
	"net"
	"strings"

	"gameclustering.com/internal/persistence"
)

type Endpoint struct {
	tcpEndpoint string
}

func handleClient(client net.Conn) {
	defer client.Close()
	buf := make([]byte, 8)
	n, err := client.Read(buf)
	if err != nil {
		//break
		return
	}
	buffer := persistence.BufferProxy{}
	buffer.NewProxy(100)
	buffer.Write(buf[:])
	buffer.Flip()
	fmt.Printf("HS %d : %d : %d\n", n, buffer.ReadInt32(), buffer.ReadInt32())
}

func (s *Endpoint) Open() error {
	parts := strings.Split(s.tcpEndpoint, "://")
	fmt.Printf("Endpoint %s %s\n", parts[0], parts[1])
	listener, err := net.Listen(parts[0], parts[1])
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		client, er := listener.Accept()
		if er != nil {
			fmt.Printf("Error :%s\n", er.Error())
			break
		}
		go handleClient(client)
	}
	return nil
}

func (s *Endpoint) Publish(event Event) error {
	event.Streaming(Chunk{true, []byte("hellop")})
	event.Streaming(Chunk{true, []byte("hellop")})
	event.Streaming(Chunk{false, []byte("hellop")})
	return nil
}
