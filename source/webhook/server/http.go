package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"webhook/config"
	"webhook/server/handler"
)

var (
	app *fiber.App
)

func Start() {

	app = fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
		},
	)

	setup()

	err := app.Listen(config.Config.Server.Addr)
	if err != nil {
		panic(err)
	}

}

func Stop() {}

func setup() {

	app.Use(recover.New())
	app.Name("webhook").Post("/webhook", handler.WebhookHandler)

}
