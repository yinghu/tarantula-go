package event

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/core"
)

type sampleFactory struct {
}

func (s *sampleFactory) Create(classId int) Event {
	return &sampleEvent{name: "sample", topic: false}
}

type sampleEvent struct {
	name     string
	topic    bool
	//listener chan Chunk
	core.PersistentableObj
}

func (s *sampleEvent) OnTopic() bool {
	return s.topic
}

func (s *sampleEvent) Streaming(c Chunk) {
	fmt.Printf("REV : %s\n", string(c.Data))
	//s.listener <- c
}

func TestEndpoint(t *testing.T) {
	buf := []byte("hello event server")
	fmt.Printf("Part : %s\n", buf[:4]) //hell
	tcp := Endpoint{TcpEndpoint: "tcp://192.168.1.4:5000", Factory: &sampleFactory{}}
	//se := sampleEvent{name: "", topic: true, listener: make(chan Chunk)}
	//go tcp.Publish(&se)
	//for ch := range se.listener {
	//fmt.Printf("F %s\n", ch.Data)
	//if !ch.Remaining {
	//break
	//}
	//}
	tcp.Open()
}
