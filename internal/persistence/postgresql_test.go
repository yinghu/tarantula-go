package persistence

import (
	"fmt"
	"testing"
)

func TestStart(t *testing.T) {
	err := Start()
	if err != nil {
		t.Errorf("SQL error %s", err.Error())
	}
}

func TestPool(t *testing.T) {
	err := Pool()
	if err != nil {
		t.Errorf("SQL error %s", err.Error())
	}
	fmt.Printf("%s\n",Message)
}
