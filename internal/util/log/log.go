package log

import (
	"fmt"

	"go.uber.org/zap"
)

func GetLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.DisableCaller = true
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.TimeKey = ""

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("failed getting logger: %s", err))
	}
	return logger
}
