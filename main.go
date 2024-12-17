package main

import (
	"github.com/adrieljss/golighter/api"
	"github.com/adrieljss/golighter/platform"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func main() {
	app := platform.App()
	fiberApp := fiber.New(platform.InitFiberConfig())
	env := app.Env

	api.SetupRoutes(fiberApp, &app)

	if err := fiberApp.Listen(env.ServerAddress); err != nil {
		log.Fatalf("server is not running: %v", err)
	}
}
