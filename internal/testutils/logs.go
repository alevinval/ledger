package testutils

import (
	"github.com/alevinval/ledger/internal/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func CaptureLogs(level zapcore.Level) (logs *observer.ObservedLogs, restore func()) {
	originalLogger := log.GetLogger().GetZapLogger()
	restoreOriginalLogger := func() {
		log.SetZapLogger(originalLogger)
	}

	capturedCore, logs := observer.New(level)
	capturedLogger := zap.New(zapcore.NewTee(originalLogger.Core(), capturedCore))
	log.SetZapLogger(capturedLogger)
	return logs, restoreOriginalLogger
}
