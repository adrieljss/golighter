package v1_users

import (
	"github.com/adrieljss/golighter/middlewares"
	"github.com/adrieljss/golighter/models"
	"github.com/adrieljss/golighter/platform"
	"github.com/adrieljss/golighter/utils"
	"github.com/gofiber/fiber/v3"
)

func SetupUsersRoutes(router fiber.Router, app *platform.Application) {
	users := router.Group("/users")

	userUsecase := NewUserUsecase(app)
	users.Get("/@me", GetMe(userUsecase), middlewares.AuthMiddleware(app, models.PermissionNone))
}

func GetMe(userUsecase *UserUsecase) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		userClaims := ctx.Locals("user").(*utils.Claims)
		user, err := userUsecase.UserRepo.GetByUID(ctx.Context(), userClaims.UID)
		if err != nil {
			return err
		}
		return ctx.Status(fiber.StatusOK).JSON(user)
	}
}

func GetUser(userUsecase *UserUsecase) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		uid := ctx.Params("uid")
		user, err := userUsecase.UserRepo.GetByUID(ctx.Context(), uid)
		if err != nil {
			return err
		}
		return ctx.Status(fiber.StatusOK).JSON(user)
	}
}

func UpdateUser(userUsecase *UserUsecase) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		uid := ctx.Params("uid")
		user, err := userUsecase.UserRepo.GetByUID(ctx.Context(), uid)
		if err != nil {
			return err
		}
		return ctx.Status(fiber.StatusOK).JSON(user)
	}
}