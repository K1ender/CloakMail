package main

import (
	"CloakMail/internal/config"
	"CloakMail/internal/logger"
	"CloakMail/internal/router"
)

func main() {
	config := config.MustInit()

	logger := logger.Init(config)
	logger.Info("Starting CloakMail...")

	router.NewRouter(config, logger)
}
