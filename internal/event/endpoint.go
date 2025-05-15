package event

import (
	"fmt"
	"net"
	"strings"
)

type Endpoint struct {
	tcpEndpoint string
}

func handleClient(client net.Conn) {
	defer client.Close()
	buf := make([]byte,1024)
	//for{
		n,err := client.Read(buf)
		if err!= nil{
			return
		}
		client.Write(buf[:n])
		//break
	//}
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
		handleClient(client)
	}
	return nil
}
