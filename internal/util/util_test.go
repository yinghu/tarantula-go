package util

import (
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
}
