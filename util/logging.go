package util

import (
	"context"
	"sync"

	zapx "go.strv.io/logging/zap"
	netx "go.strv.io/net"
	loggerx "go.strv.io/net/logger"

	"go.uber.org/zap"
)

const (
	requestIDFieldName = "request_id"
)

var (
	// once ensures that no one will override server log level.
	once = &sync.Once{}

	// serverLogLevel is by default info log level.
	serverLogLevel = zap.NewAtomicLevel()
)

// SetServerLogLevel sets log level that will be used for all new server loggers.
// Only first call of this function is valid. Other calls are ignored to prevent unwanted behavior.
func SetServerLogLevel(l zap.AtomicLevel) {
	once.Do(func() {
		serverLogLevel = l
	})
}

// ServerLogger is a logger that is used by go.strv.io/net/http.Server to log errors and debug messages.
type ServerLogger struct {
	*zap.Logger
}

// NewServerLogger returns new instance of ServerLogger with given caller value.
func NewServerLogger(caller string) ServerLogger {
	l := zapx.MustCreateLogger(zapx.Config{
		Level:             serverLogLevel,
		DisableStacktrace: true,
		DisableCaller:     true,
	})
	l = l.With(
		zap.String("caller", caller),
	)
	return ServerLogger{l}
}

// With is a wrapper around zap.With using zap.Any as a field type.
func (l ServerLogger) With(fields ...loggerx.Field) loggerx.ServerLogger {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return ServerLogger{Logger: l.Logger.With(zapFields...)}
}

// Debug logs debug message.
func (l ServerLogger) Debug(msg string) {
	l.Logger.Debug(msg)
}

// Info logs info message.
func (l ServerLogger) Info(msg string) {
	l.Logger.Info(msg)
}

// Warn logs warning message.
func (l ServerLogger) Warn(msg string) {
	l.Logger.Warn(msg)
}

// Error logs error message.
func (l ServerLogger) Error(msg string, err error) {
	l.Logger.Error(msg, zap.Error(err))
}

// WithCtx returns logger with fields extracted from context.
func WithCtx(ctx context.Context, l *zap.Logger) *zap.Logger {
	return l.With(
		zap.String(requestIDFieldName, netx.RequestIDFromCtx(ctx)),
	)
}
