package app

import (
	"io"
	"log/slog"
)

// NoopLogger returns a no-op logger which discards all logs
func NoopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
