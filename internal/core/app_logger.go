package core

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

var (
	AppLog zerolog.Logger
)

func CreateAppLog(dir string, truncated bool) {
	fmt.Printf("Creating app log %s\n", dir)
	err := os.MkdirAll(dir+"/log", 0755)
	if err != nil {
		AppLog = zerolog.New(os.Stdout)
		return
	}
	opt := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	if truncated {
		opt = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	file, err := os.OpenFile(dir+"/log/tarantula.log", opt, 0644)
	if err != nil {
		AppLog = zerolog.New(os.Stdout)
		return
	}
	AppLog = zerolog.New(file)
	AppLog.Info().Msg("Initialized app log")
}

func CreateTestLog() {
	AppLog = zerolog.New(os.Stdout)
	AppLog.Info().Msg("Initialized app log")
}
