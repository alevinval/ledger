//go:build !debug

package log

import (
	"fmt"

	"go.uber.org/zap"
)

func newZapLogger(level zap.AtomicLevel) *zap.Logger {
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.Level = level
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.TimeKey = ""

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("failed getting logger: %s", err))
	}
	return logger
}
