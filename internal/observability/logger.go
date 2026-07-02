package observability

import (
	"log/slog"
	"os"
)

// NewLogger cria um logger estruturado que escreve em JSON no stdout.
func NewLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return slog.New(handler)
}
