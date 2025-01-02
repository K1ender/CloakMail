package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func NewHealthHandler(router fiber.Router, logger *logrus.Logger) {
	logger.Info("Registering health check route")
	router.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message":   "ok",
			"timestamp": time.Now(),
		})
	})
}
