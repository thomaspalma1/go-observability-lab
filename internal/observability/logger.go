package observability

import (
	"log/slog"
	"os"
)

// NewLogger creates a structured logger that writes JSON logs to stdout.
func NewLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return slog.New(handler)
}
