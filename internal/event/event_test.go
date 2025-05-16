package event

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/persistence"
)

type sampleEvent struct {
	name     string
	topic    bool
	listener chan Chunk
	persistence.PersistentableObj
}

func (s *sampleEvent) Topic() bool {
	return s.topic
}

func (s *sampleEvent) Streaming(c Chunk) {
	s.listener <- c
}

func TestEndpoint(t *testing.T) {
	buf := []byte("hello event server")
	fmt.Printf("Part : %s\n", buf[:4]) //hell
	tcp := Endpoint{tcpEndpoint: "tcp://192.168.1.4:5000"}
	se := sampleEvent{name: "", topic: true, listener: make(chan Chunk)}
	go tcp.Publish(&se)
	for ch := range se.listener {
		fmt.Printf("F %s\n", ch.Data)
		if !ch.Remaining {
			break
		}
	}
}
