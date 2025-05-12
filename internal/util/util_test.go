package util

import (
	"fmt"
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
	nodes := []string{"a01", "a02", "a03", "a04", "a05", "a06", "a07", "a08", "a09", "a10", "a11", "a12"}

	for n := range nodes {
		fmt.Printf("Patition %d %s\n", n%7, nodes[n])
	}
}
