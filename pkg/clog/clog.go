package clog

import (
	"context"
	"log/slog"
)

type key struct{}

var loggerKey key

func WithContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func FromContext(ctx context.Context) *slog.Logger {
	l, exists := ctx.Value(loggerKey).(*slog.Logger)
	if !exists || l == nil {
		panic("clog logger not setup")
	}

	return l
}
