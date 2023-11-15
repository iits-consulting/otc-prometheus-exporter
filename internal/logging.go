package internal

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ILogger interface {
	Info(msg string, tags ...interface{})
	Debug(msg string, tags ...interface{})
	Warn(msg string, tags ...interface{})
	Error(msg string, tags ...interface{})
	Panic(msg string, tags ...interface{})
	Sync() error
}

type Logger struct {
	internalLogger *zap.Logger
}

func NewLogger(logLevel string) ILogger {
	var level zapcore.Level
	switch logLevel {
	case "DEBUG":
		level = zapcore.DebugLevel
	case "INFO":
		level = zapcore.InfoLevel
	case "WARN":
		level = zapcore.WarnLevel
	case "ERROR":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	config := zap.NewProductionConfig()

	// Set level to our chosen level above
	config.Level = zap.NewAtomicLevelAt(level)

	internalLogger, _ := config.Build()
	return &Logger{internalLogger: internalLogger}
}

func (l *Logger) Info(msg string, tags ...interface{}) {
	l.internalLogger.Info(msg, tagsToZapFields(tags...)...)
}

func (l *Logger) Debug(msg string, tags ...interface{}) {
	l.internalLogger.Debug(msg, tagsToZapFields(tags...)...)
}

func (l *Logger) Warn(msg string, tags ...interface{}) {
	l.internalLogger.Warn(msg, tagsToZapFields(tags...)...)
}

func (l *Logger) Error(msg string, tags ...interface{}) {
	l.internalLogger.Error(msg, tagsToZapFields(tags...)...)
}

func (l *Logger) Panic(msg string, tags ...interface{}) {
	l.internalLogger.Panic(msg, tagsToZapFields(tags...)...)
}

func (l *Logger) Sync() error {
	return l.internalLogger.Sync()
}

func tagsToZapFields(tags ...interface{}) (field []zap.Field) {
	var zapFields []zap.Field

	// Iterate over tags and create appropriate fields so that we have json output

	for i := 0; i < len(tags); i += 2 {
		key, didCast := tags[i].(string)
		value := tags[i+1]

		if didCast {
			zapFields = append(zapFields, zap.Any(key, value))
		} else {
			panic("( logging.go -> tagsToZapFields() ) Tag key must be string to convert to zap field.")
		}
	}

	return zapFields
}
