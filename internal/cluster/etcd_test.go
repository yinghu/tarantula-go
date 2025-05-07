package cluster

import (
	//"fmt"
	"testing"
	"time"
)

func TestLink(t *testing.T) {
	c := Cluster{Group:"tarantula",Endpoints: []string{"192.168.1.7:2379"}, Timeout: 5 * time.Second}
	err := c.Watch()
	if err != nil {
		t.Errorf("Service error %s", err.Error())
	}
}
