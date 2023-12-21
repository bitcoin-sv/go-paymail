package logging

import (
	"github.com/rs/zerolog"
	"go.elastic.co/ecszerolog"
	"os"
)

// GetDefaultLogger generates and returns a default logger instance.
func GetDefaultLogger() *zerolog.Logger {
	logger := ecszerolog.New(os.Stdout, ecszerolog.Level(zerolog.DebugLevel)).
		With().
		Timestamp().
		Caller().
		Str("application", "go-paymail").
		Logger()

	return &logger
}
