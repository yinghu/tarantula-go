package cluster

import (
	//"fmt"
	"testing"
	"time"
)

func TestCluster(t *testing.T) {
	c := NewEtc("tarantula",[]string{"192.168.1.7:2379"},Node{Name: "a01", HttpEndpoint: "http://192.168.1.11:8080", TcpEndpoint: "tcp://192.168.1.11:5000"})

	tk := time.NewTimer(30 * time.Second)
	go func() {
		<-tk.C
		c.Quit <- true
	}()
	err := c.Join()
	if err != nil {
		t.Errorf("Service error %s", err.Error())
	}
}
