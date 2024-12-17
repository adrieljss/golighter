package api

import (
	api_v1 "github.com/adrieljss/golighter/api/v1"
	"github.com/adrieljss/golighter/platform"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func SetupRoutes(fiberApp *fiber.App, app *platform.Application) {
	api := fiberApp.Group("/api")
	api.Use("/", AcceptJson())
	api.Use("/", logger.New())

	api_v1.SetupApiRoutes(api, app)
}

func AcceptJson() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		ctx.Accepts("application/json")
		return ctx.Next()
	}
}
