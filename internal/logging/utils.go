// Package logging provides helpers for working with loggers down call chains
package logging

import (
	"context"
	"log/slog"
)

type ctxKey string

var loggerCtxKey ctxKey = "logger"

func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerCtxKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
