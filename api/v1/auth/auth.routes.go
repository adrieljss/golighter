package v1_auth

import (
	"errors"

	"github.com/adrieljss/golighter/platform"
	"github.com/adrieljss/golighter/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func SetupAuthRoutes(router fiber.Router, app *platform.Application) {
	auth := router.Group("/auth")

	authUsecase := NewAuthUsecase(app)
	auth.Post("/register", RegisterUser(authUsecase))
	auth.Post("/login", LoginUser(authUsecase))
	auth.Post("/refresh", RefreshToken(authUsecase))
}

func RegisterUser(authUsecase *AuthUsecase) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var user UserRegister
		err := ctx.Bind().Body(&user)
		if err != nil {
			return err
		}

		response, err := authUsecase.CreateUser(ctx.Context(), &user)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(response)
	}
}

func LoginUser(authUsecase *AuthUsecase) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var user UserLogin
		err := ctx.Bind().Body(&user)
		if err != nil {
			return err
		}

		response, err := authUsecase.LoginUser(ctx.Context(), &user)
		if err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || errors.Is(err, pgx.ErrNoRows) {
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":    "invalid credentials",
					"metadata": utils.NewMetadata().Set("email", "invalid credentials").Set("username", "invalid credentials"),
				})
			}
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(response)
	}
}

func RefreshToken(authUsecase *AuthUsecase) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var user UserRefreshToken
		err := ctx.Bind().Body(&user)
		if err != nil {
			return err
		}

		response, err := authUsecase.RefreshAccessToken(ctx.Context(), user.RefreshToken)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(AccTokenResponse{
			AccessToken: response,
		})
	}
}
