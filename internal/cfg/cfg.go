package cfg

import (
	"context"
	"efmob/logger"
	"sync"

	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

var (
	once   sync.Once
	config Config
)

type Config struct {
	DbDSN string `env:"DATABASE_URI"`
	Port  string `env:"PORT"`
}

func GetConfig() *Config {
	once.Do(func() {
		if err := envconfig.Process(context.Background(), &config); err != nil {
			logger.Log().Fatal("Failed to load config from env", zap.Error(err))
		}
	})

	return &config
}
