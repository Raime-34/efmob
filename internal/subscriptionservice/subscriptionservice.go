package subscriptionservice

import (
	"context"
	"efmob/internal/cfg"
	"efmob/internal/migrate"
	"efmob/internal/repositories"
	"efmob/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type SubService struct {
	repo Repositories
}

func InitService() *SubService {
	dbConfig, err := pgxpool.ParseConfig(cfg.GetConfig().DbDSN)
	if err != nil {
		logger.Log().Fatal("parse db dsn", zap.Error(err))
	}

	connPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		logger.Log().Fatal("pgx pool", zap.Error(err))
	}

	migrate.MakeMigration(connPool)

	return &SubService{
		repo: repositories.NewRepo(connPool),
	}
}

type Repositories interface{}
