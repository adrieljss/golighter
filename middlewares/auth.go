package middlewares

import (
	"strings"

	"github.com/adrieljss/golighter/platform"
	"github.com/adrieljss/golighter/utils"
	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware(app *platform.Application) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "unauthorized",
			})
		}
		// can start with "Bearer " or none
		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateJWT(token, app.Env.JWTSecretAccessToken)
		if err != nil {
			// abort request
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "unauthorized",
			})
		}

		c.Locals("user", claims)

		return c.Next()
	}
}
