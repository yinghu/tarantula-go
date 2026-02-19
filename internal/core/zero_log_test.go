package core

import (
	"testing"
)


func TestZeroLog(t *testing.T) {
	CreateTestLog()
	AppLog.Debug().Msgf("test %s","a test")
	AppLog.Printf("test %s","mess")

}
