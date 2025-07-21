package logger

import (
	"context"
	"log/slog"
	"os"
)

const (
	RequestPath       = "request_path"
	RequestMethod     = "request_method"
	RequestRemoteAddr = "request_remote_address"

	APIMethod  = "api_method"
	TraceID    = "trace_id"
	StatusCode = "status_code"
	Error      = "error"

	UserID    = "user_id"
	UserLogin = "user_login"

	HandlerStartedEvent    = "handler started"
	HandlerCompletedEvent  = "handler completed"
	HandlerErrorEvent      = "handler error"
	GetTokenAmountEvent    = "get client token amount"
	RemainTokenAmountEvent = "remain client token amount"
)

type keyType int

const loggerKey = keyType(0)

func FromContext(ctx context.Context) *slog.Logger {
	v := ctx.Value(loggerKey)
	if v == nil {
		return slog.Default()
	}

	logger := v.(*slog.Logger)
	return logger
}

func ContextWithSlogLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func AddValuesToContext(ctx context.Context, args ...any) context.Context {
	l := FromContext(ctx)
	l = l.With(args...)
	return ContextWithSlogLogger(ctx, l)
}

func InfoAddValues(ctx context.Context, message string, args ...any) context.Context {
	l := FromContext(ctx)
	l = l.With(args...)
	l.Info(message)
	return ContextWithSlogLogger(ctx, l)
}

func DebagAddValues(ctx context.Context, message string, args ...any) context.Context {
	l := FromContext(ctx)
	l = l.With(args...)
	l.Debug(message)
	return ContextWithSlogLogger(ctx, l)
}

func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

func InitLogging() {
	handler := slog.Handler(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
	))

	slog.SetDefault(slog.New(handler))
}
