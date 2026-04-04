package bootstrap

import (
	"os"

	"gameclustering.com/internal/core"
	"github.com/rs/zerolog"
)

func CreateAppLog(dir string, truncated bool, standAlone bool) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
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
	//zerolog.MultiLevelWriter()
	core.AppLog = zerolog.New(file).With().Timestamp().Logger().With().Caller().Logger()
	core.AppLog.Info().Msg("Initialized app log")
}

func CreateTestLog() {
	core.AppLog = zerolog.New(os.Stderr).With().Timestamp().Logger().With().Caller().Logger()
	core.AppLog.Info().Msg("Initialized test app log")
}
