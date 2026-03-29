package core

import (
	"os"

	"github.com/rs/zerolog"
)

var (
	AppLog zerolog.Logger
)

func CreateAppLog(dir string, truncated bool, standAlone bool) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if standAlone {
		CreateTestLog()
		return
	}
	err := os.MkdirAll(dir+"/log", 0755)
	if err != nil {
		CreateTestLog()
		return
	}
	opt := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	if truncated {
		opt = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	file, err := os.OpenFile(dir+"/log/tarantula.log", opt, 0644)
	if err != nil {
		CreateTestLog()
		return
	}
	AppLog = zerolog.New(file).With().Timestamp().Logger().With().Caller().Logger()
	AppLog.Info().Msg("Initialized app log")
}

func CreateTestLog() {
	AppLog = zerolog.New(os.Stderr).With().Timestamp().Logger().With().Caller().Logger()
	AppLog.Info().Msg("Initialized test app log")
}
