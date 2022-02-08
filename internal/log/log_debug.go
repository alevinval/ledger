//go:build debug

package log

import "fmt"
import "go.uber.org/zap"

func newZapLogger(level zap.AtomicLevel) *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.DisableStacktrace = true

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("failed getting logger: %s", err))
	}
	return logger
}
