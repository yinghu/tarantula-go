package bootstrap

import (
	"io"
	"os"

	"gameclustering.com/internal/core"
	"github.com/rs/zerolog"
)

func CreateAppLog(dir string, truncated bool, standAlone bool, out io.Writer) {
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
	core.AppLog = zerolog.New(zerolog.MultiLevelWriter(file, out)).With().Timestamp().Caller().Logger()
	core.AppLog.Info().Msg("Initialized app log")
}

func CreateTestLog() {
	core.AppLog = zerolog.New(os.Stderr).With().Timestamp().Logger().With().Caller().Logger()
	core.AppLog.Info().Msg("Initialized test app log")
}
