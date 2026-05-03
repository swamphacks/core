package logger

import (
	"os"
	"runtime/debug"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New() zerolog.Logger {
	buildInfo, _ := debug.ReadBuildInfo()

	logger := zerolog.New(os.Stderr).Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()

	// set global logger too for now
	log.Logger = logger

	return logger
}
