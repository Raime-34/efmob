// Package main — точка входа сервиса подписок.
//
// @title Subscription Service API
// @version 1.0
// @description REST API подписок: создание, изменение, удаление, список, сумма по фильтру.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
package main

import (
	"efmob/internal/cfg"
	"efmob/internal/subscriptionservice"
	"efmob/logger"
	"os"
	"os/signal"
	"syscall"

)

func main() {
	logger.Log().Info("starting subscriptionservice")

	cfg.GetConfig()
	subscriptionservice.InitService()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log().Info("shutdown complete")
}
