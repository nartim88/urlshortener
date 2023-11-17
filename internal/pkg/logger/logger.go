package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var Log = zerolog.Nop()

func Init(level string) error {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger := zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339Nano,
		},
	)

	logLevel, err := getLogLevel(level)
	if err != nil {
		return err
	}

	logger = logger.Level(*logLevel).With().
		Timestamp().
		Caller().
		Int("pid", os.Getgid()).
		Logger()

	Log = logger

	return nil
}

func getLogLevel(level string) (*zerolog.Level, error) {
	logLevels := map[string]zerolog.Level{
		"trace":   zerolog.TraceLevel,
		"debug":   zerolog.DebugLevel,
		"info":    zerolog.InfoLevel,
		"warning": zerolog.WarnLevel,
		"error":   zerolog.ErrorLevel,
		"fatal":   zerolog.FatalLevel,
		"panic":   zerolog.PanicLevel,
	}

	zLogLevel, ok := logLevels[level]
	if !ok {
		return nil, LogLevelError{"the given log level doesn't supported by logger", level}
	}

	return &zLogLevel, nil
}
