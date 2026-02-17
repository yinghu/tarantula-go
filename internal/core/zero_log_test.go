package core

import (
	"testing"
)


func TestZeroLog(t *testing.T) {
	CreateTestLog()
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	//log.Logger.Printf("test %d",100)
	//log.Info().Str("key","value").Float64("money",100).Msg("tranaction");
	//logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	//logger.Printf("TEST %s","mell")
	//logger.Debug().Msg("another mee")
	AppLog.Printf("test %s","mess")

}
