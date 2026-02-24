package core

import (
	"testing"
)

func TestZeroLog(t *testing.T) {
	CreateTestLog()
	AppLog.Debug().Msgf("test %s", "a test")
	AppLog.Printf("test %s", "mess")
	var data [][]byte
	data = append(data, []byte("key1"))
	data = append(data, []byte("key2"))
	//data = append(data, []byte("key3"))
	//data = append(data, []byte("key4"))
	sz := min(3,len(data))
	cy := make([][]byte, 0, sz)
	for i:= range sz{
		cy = append(cy,data[i])
	}
	data = data[sz:]
	//AppLog.Debug().Msgf("copied %d",n)
	for _, d := range data{
		AppLog.Debug().Msgf("data %s", string(d))
	}

	for _, d := range cy{
		AppLog.Debug().Msgf("copy %s", string(d))
	}
}
