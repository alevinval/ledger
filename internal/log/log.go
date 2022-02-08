package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zapLogLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	instance    = &Logger{zapLogger: newZapLogger(zapLogLevel)}
)

type Logger struct {
	zapLogger *zap.Logger
}

func GetLogger() *Logger {
	return instance
}

func SetZapLogger(logger *zap.Logger) {
	instance.zapLogger = logger
}

func (l *Logger) GetZapLogger() *zap.Logger {
	return l.zapLogger
}

func (l *Logger) SetLevel(level zapcore.Level) {
	zapLogLevel.SetLevel(level)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zapLogger.Fatal(msg, fields...)
}
