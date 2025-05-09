package cluster

import (
	//"fmt"
	"testing"
	"time"
)

func TestLink(t *testing.T) {
	c := Cluster{Group: "tarantula", EtcdEndpoints: []string{"192.168.1.7:2379"}, Timeout: 5 * time.Second, Local: Node{Name: "a01", HttpEndpoint: "http://192.168.1.11:8080", TcpEndpoint: "tcp://192.168.1.11:5000"}}
	err := c.Join()
	if err != nil {
		t.Errorf("Service error %s", err.Error())
	}
}
