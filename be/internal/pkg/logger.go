package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey string

const (
	// RequestIDKey is the key used to store the requestID in the context
	RequestIDKey ctxKey = "request_id"
	severityKey         = "severity"
)

var (
	// Global logger instance
	log *zap.Logger
)

// Init initializes the global logger
// If enableDebug is true, DEBUG level logs will be included, otherwise minimal level will be INFO
func Init(enableDebug bool) {
	cfg := zap.NewProductionConfig()

	// Set log level based on enableDebug parameter
	if enableDebug {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Configure for GCP Logging compatibility
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.LevelKey = severityKey

	// Map Zap levels to GCP severity levels
	cfg.EncoderConfig.EncodeLevel = zapcore.LevelEncoder(func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString("DEBUG")
		case zapcore.InfoLevel:
			enc.AppendString("INFO")
		case zapcore.WarnLevel:
			enc.AppendString("WARNING")
		case zapcore.ErrorLevel:
			enc.AppendString("ERROR")
		case zapcore.DPanicLevel, zapcore.PanicLevel:
			enc.AppendString("CRITICAL")
		case zapcore.FatalLevel:
			enc.AppendString("EMERGENCY")
		default:
			enc.AppendString("DEFAULT")
		}
	})

	var err error
	log, err = cfg.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}

// WithRequestID creates a logger with the requestID field
func WithRequestID(ctx context.Context) *zap.Logger {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok || requestID == "" {
		return log
	}
	return log.With(zap.String("request_id", requestID))
}

// ContextWithTraceID returns a new context with the traceID value
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	// Create a background context if nil is provided
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetTraceID gets the traceID from context
func GetTraceID(ctx context.Context) string {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return ""
	}
	return requestID
}

// Debug logs a debug message
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	WithRequestID(ctx).Debug(msg, fields...)
}

// Info logs an info message
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	WithRequestID(ctx).Info(msg, fields...)
}

// Warn logs a warning message
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	WithRequestID(ctx).Warn(msg, fields...)
}

// Error logs an error message
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	WithRequestID(ctx).Error(msg, fields...)
}

// Fatal logs a fatal message and then calls os.Exit(1)
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	WithRequestID(ctx).Fatal(msg, fields...)
}

// With creates a child logger and adds structured context to it
func With(ctx context.Context, fields ...zap.Field) *zap.Logger {
	return WithRequestID(ctx).With(fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return log.Sync()
}
