package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "DevPrep",
		AppName:       "Dev Prep app  v1.0.1",
	})

	app.Listen(":3000")
}
