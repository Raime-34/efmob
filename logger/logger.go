package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	once   sync.Once
	logger *zap.Logger
)

func Log() *zap.Logger {
	once.Do(func() {
		logger, _ = zap.NewProduction()
	})

	return logger
}
