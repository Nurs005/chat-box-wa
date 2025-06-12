package infrastructure

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Logger = zerolog.Logger

func NewLogger() Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(i.(string))
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	return logger
}
