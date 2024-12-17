package v1_health

import (
	"github.com/gofiber/fiber/v3"
)

func HealthCheck(ctx fiber.Ctx) error {
	return ctx.SendString("OK")
}

func SetupHealthRoutes(router fiber.Router) {
	health := router.Group("/health")
	health.Get("/", HealthCheck)
}
