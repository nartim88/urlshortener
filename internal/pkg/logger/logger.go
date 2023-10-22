package logger

import (
	"os"

	"github.com/rs/zerolog"
)

const defLogLevel = zerolog.InfoLevel

var Log zerolog.Logger = zerolog.New(os.Stdout).Level(defLogLevel).With().Timestamp().Logger()

func Init(level string) error {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logLevel, err := getLogLevel(level)
	if err != nil {
		return err
	}
	logger.Level(*logLevel)
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
		return nil, LogLevelError{"the given log level doesn't supported by logger"}
	}

	return &zLogLevel, nil
}
