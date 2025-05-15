package event

import (
	"fmt"
	"testing"
)

func TestEndpoint(t *testing.T) {
	buf := []byte("hello event server")
	fmt.Printf("Part : %s\n", buf[:4]) //hell
	tcp := Endpoint{tcpEndpoint: "tcp://192.168.1.4:5000"}
	err := tcp.Open()
	if err != nil {
		t.Errorf("Error %s\n", err.Error())
	}
}
