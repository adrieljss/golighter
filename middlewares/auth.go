package middlewares

import (
	"strings"

	"github.com/adrieljss/golighter/models"
	"github.com/adrieljss/golighter/platform"
	"github.com/adrieljss/golighter/utils"
	"github.com/gofiber/fiber/v3"
)

// AuthMiddleware is a middleware that checks if the user is authenticated and has the required permissions
// requiredPerms is a bitmask of the required permissions, the user must have all the bitmask permissions, else it will be forbidden.
func AuthMiddleware(app *platform.Application, requiredPerms models.Permission) fiber.Handler {
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

		if claims.PermissionFlags&requiredPerms != requiredPerms {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "forbidden",
			})
		}

		c.Locals("user", claims)

		return c.Next()
	}
}
