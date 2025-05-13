package util

import (
	"fmt"
	"slices"
	"testing"
)

func TestPassword(t *testing.T) {
	h, err := Hash("password")
	if err != nil {
		t.Errorf("failed %s\n", err.Error())
	}
	er := Match("password", h)
	if er != nil {
		t.Errorf("failed %s\n", er.Error())
	}
	nodes := []string{"a11", "a02", "a03", "a04", "a05", "a06", "a07", "a08", "a09", "a10", "a01", "a12"}
	slices.Sort(nodes)
	mp := make(map[int]string)
	sz := len(nodes)
	for p:= range 31{
		i := p%sz
		fmt.Printf("Partition %d %s %d\n",i,nodes[i],p)
		mp[p]=nodes[i]
	}
}
