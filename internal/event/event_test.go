package event

import (
	"testing"
)

func TestEndpoint(t *testing.T) {
	tcp := Endpoint{tcpEndpoint: "tcp://192.168.1.4:5000"}
	err := tcp.Open()
	if err != nil {
		t.Errorf("Error %s\n", err.Error())
	}
}
