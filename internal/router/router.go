package router

import (
	"CloakMail/internal/config"
	"CloakMail/internal/handlers"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func NewRouter(cfg *config.Config, logger *logrus.Logger) {
	logger.Info("Starting router...")

	app := fiber.New(fiber.Config{})

	v1 := app.Group("/api/v1")

	handlers.NewHealthHandler(
		v1,
		logger.WithField("route", "health_check").Logger,
	)

	logger.Fatal(app.Listen(
		fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
	))
}
