package main

import (
	"github.com/adrieljss/golighter/api"
	"github.com/adrieljss/golighter/platform"
	"github.com/gofiber/fiber/v3/log"
)

func main() {
	app := platform.App()
	api.SetupRoutes(&app)

	if err := app.FiberApp.Listen(app.Env.ServerAddress); err != nil {
		log.Fatalf("server is not running: %v", err)
	}
}
