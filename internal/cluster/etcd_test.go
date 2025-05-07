package cluster

import (
	//"fmt"
	"testing"
)

func TestLink(t *testing.T) {
	err := Link()
	if err!=nil{
		t.Errorf("Service error %s", err.Error())
	}
}
