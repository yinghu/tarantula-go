package util

import (
	"testing"
)

func TestPassword(t *testing.T) {
	h, err := HashPassword("password")
	if err != nil {
		t.Errorf("failed %s\n", err.Error())
	}
	er := ValidatePassword("password", h)
	if er != nil {
		t.Errorf("failed %s\n", er.Error())
	}
}

func TestPartition(t *testing.T) {
	p := Partition([]byte("hellp"),5)
	if p>5 {
		t.Errorf("falied partition %d\n",p)
	}
}