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
	log := logger.Log()
	log.Info("starting subscriptionservice")

	cfg.GetConfig()
	subscriptionservice.InitService()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutdown complete")
}
