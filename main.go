package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello World!")
	})

	err := app.Listen("localhost:3000")
	if err != nil {
		panic(err)
	}
}
