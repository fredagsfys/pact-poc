package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/term"
)

// Config contains logger configuration.
type Config struct {
	Level string `yaml:"level"`
}

// Logger allows to write logs.
type Logger struct {
	logger zerolog.Logger
}

// New returns new logger.
func New(c *Config) (*Logger, error) {
	if c == nil {
		c = &Config{}
	}

	level := zerolog.InfoLevel

	if c.Level != "" {
		var err error
		level, err = zerolog.ParseLevel(c.Level)
		if err != nil {
			return nil, fmt.Errorf("unknown level '%s'", c.Level)
		}
	}

	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"

	file := os.Stdout

	logger := zerolog.New(file).Level(level)

	logger = logger.With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + 1).
		Logger()

	if term.IsTerminal(int(file.Fd())) {
		logger = logger.Output(zerolog.ConsoleWriter{
			Out:        file,
			TimeFormat: time.RFC3339,
		})
	}

	return &Logger{
		logger: logger,
	}, nil
}

// Debug prints debug logs
func (l *Logger) Debug(format string, v ...interface{}) {
	l.logger.Debug().Msgf(format, v...)
}

// Info prints info logs
func (l *Logger) Info(format string, v ...interface{}) {
	l.logger.Info().Msgf(format, v...)
}

// Warn prints warn logs
func (l *Logger) Warn(format string, v ...interface{}) {
	l.logger.Warn().Msgf(format, v...)
}

// Error prints error logs
func (l *Logger) Error(format string, v ...interface{}) {
	l.logger.Error().Msgf(format, v...)
}
