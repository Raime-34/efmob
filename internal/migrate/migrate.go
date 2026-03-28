package migrate

import (
	"efmob/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func MakeMigration(conn *pgxpool.Pool) {
	db := stdlib.OpenDBFromPool(conn)
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Log().Fatal("db ping", zap.Error(err))
	}

	if err := goose.Up(db, "/app//migrations"); err != nil {
		logger.Log().Fatal("goose up", zap.Error(err))
	}
}
