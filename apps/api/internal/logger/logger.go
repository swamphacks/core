package logger

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"
	_ "time/tzdata"

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

func NewCheckinLogger() zerolog.Logger {
	eastern, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}

	output := zerolog.ConsoleWriter{
		Out:           os.Stdout,
		NoColor:       true,
		TimeLocation:  eastern,
		TimeFormat:    "2006-01-02 15:04:05 MST",
		PartsOrder:    []string{"time", "type", "message"},
		FieldsExclude: []string{"type", "pid", "go_version"},
		FormatPartValueByName: func(value interface{}, name string) string {
			if name == "type" && value != nil {
				return fmt.Sprintf("[%v]", value)
			}
			return ""
		},
	}
	return zerolog.New(output).With().Timestamp().Logger()
}
