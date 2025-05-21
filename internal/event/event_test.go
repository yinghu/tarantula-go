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
func (s *sampleFactory) OnEvent(e Event){
	//&sampleEvent{name: "sample", topic: false}
}

type sampleEvent struct {
	name  string
	topic bool
	//listener chan Chunk
	core.PersistentableObj
	EventObj
}

func (s *sampleEvent) Read(buffer core.DataBuffer) error {
	a, _ := buffer.ReadInt32()
	b, _ := buffer.ReadInt32()
	fmt.Printf("Reading from buffer %d, %d\n", a, b)
	return nil
}

func (s *sampleEvent) OnTopic() bool {
	return s.topic
}

func (s *sampleEvent) streaming(c Chunk) {
	fmt.Printf("REV : %s\n", string(c.Data))
	//s.listener <- c
}

func (s *sampleEvent) Inbound(buff core.DataBuffer) {
	s.Read(buff)
	for {
		sz, err := buff.ReadInt32()
		if err != nil {
			s.streaming(Chunk{true, []byte{0}})
			break
		}
		if sz == 0 {
			s.streaming(Chunk{true, []byte{0}})
			break
		}
		pd, err := buff.Read(int(sz))
		if err != nil {
			s.streaming(Chunk{true, []byte{0}})
			break
		}
		s.streaming(Chunk{false, pd})
	}
}

func (s *sampleEvent) Outbound(buff core.DataBuffer) {
	buff.WriteInt32(100)
	buff.WriteString("Bye")
}

func TestEndpoint(t *testing.T) {
	buf := []byte("hello event server")
	fmt.Printf("Part : %s\n", buf[:4]) //hell
	tcp := Endpoint{TcpEndpoint: "tcp://192.168.1.4:5000", Service: &sampleFactory{}}
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
