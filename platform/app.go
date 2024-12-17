package platform

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/adrieljss/golighter/utils"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// for read only, it's safe to use pointers
type Application struct {
	Env    *Env          // read-only
	Db     *pgxpool.Pool // read-only
	Mailer *SMTPMailer   // read-only
}

// initialize env, db, mailer (in that order)
func App() Application {
	app := Application{}

	app.Env = NewEnv()
	app.Db = ConnectDB(initPgConfig(app.Env))
	app.Mailer = NewMailer(initMailerConfig(app.Env))

	return app
}

func initPgConfig(env *Env) PgConfig {
	return PgConfig{
		Host:     env.DBHost,
		Port:     env.DBPort,
		User:     env.DBUser,
		Password: env.DBPassword,
		DbName:   env.DBName,
	}
}

func initMailerConfig(env *Env) SmtpConfig {
	return SmtpConfig{
		Host:         env.SMTPHost,
		Port:         env.SMTPPort,
		EmailAddress: env.SMTPFrom,
		Password:     env.SMTPPassword,
	}
}

type structValidator struct {
	validate *validator.Validate
}

// Validator needs to implement the Validate method
func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func InitFiberConfig() fiber.Config {
	en := en.New()
	trans, _ := ut.New(en, en).GetTranslator("en")
	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	return fiber.Config{
		JSONEncoder:     json.Marshal,
		JSONDecoder:     json.Unmarshal,
		StructValidator: &structValidator{validate},
		ErrorHandler:    errorHandlerFunc(trans),
	}
}

func errorHandlerFunc(trans ut.Translator) fiber.ErrorHandler {
	return func(ctx fiber.Ctx, err error) error {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			metadata := utils.NewMetadata()
			for _, err := range validationErrors {
				metadata.Set(err.Tag(), err.Translate(trans))
			}
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":    "validation error",
				"metadata": metadata,
			})
		}

		code := fiber.StatusInternalServerError
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "not found",
			})
		}

		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			// TODO: Add more error codes (if necessary)
			switch pgError.Code {
			case "23505":
				code = fiber.StatusBadRequest
				column := utils.GetColumnFromConstraint(pgError)
				return ctx.Status(code).JSON(fiber.Map{
					"error":    fmt.Sprintf("%s already exists", column),
					"metadata": utils.NewMetadata().Set(column, fmt.Sprintf("%s already exists", column)),
				})
			}
		}

		log.Error(err)
		return ctx.Status(code).JSON(fiber.Map{
			"error": "unexpected error, please try again later",
		})
	}
}
