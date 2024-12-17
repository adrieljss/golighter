package v1_api

import (
	v1_auth "github.com/adrieljss/golighter/api/v1/auth"
	v1_health "github.com/adrieljss/golighter/api/v1/health"
	v1_users "github.com/adrieljss/golighter/api/v1/users"
	"github.com/adrieljss/golighter/platform"
	"github.com/gofiber/fiber/v3"
)

func SetupApiRoutes(router fiber.Router, app *platform.Application) {
	api := router.Group("/v1")

	v1_auth.SetupAuthRoutes(api, app)
	v1_health.SetupHealthRoutes(api)
	v1_users.SetupUsersRoutes(api, app)
}
