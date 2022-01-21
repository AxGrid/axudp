package axtools

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func GetLogLevel(levelString string) zerolog.Level {
	switch levelString {
	default:
		return zerolog.DebugLevel
	case "trace":
		return zerolog.TraceLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "err":
		return zerolog.ErrorLevel
	}
}

func InitLogger(level string) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05,000"}).Level(GetLogLevel(level))
	return log.Logger
}
