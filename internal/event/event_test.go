package event

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/persistence"
)

type sampleEvent struct {
	name  string
	topic bool
	mc    chan []byte
	persistence.PersistentableObj
}

func (s *sampleEvent) Topic() bool {
	return s.topic
}

func (s *sampleEvent) Send() error {
	return nil
}


func TestEndpoint(t *testing.T) {
	buf := []byte("hello event server")
	fmt.Printf("Part : %s\n", buf[:4]) //hell
	tcp := Endpoint{tcpEndpoint: "tcp://192.168.1.4:5000"}
	se := sampleEvent{name: "",topic: true,mc:make(chan []byte)}
	tcp.Publish(&se)
	//err := tcp.Open()
	//if err != nil {
	//t.Errorf("Error %s\n", err.Error())
	//}
}
