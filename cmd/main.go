package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/ycd/leaderboard/pkg/handlers"
)

func main() {
	app := fiber.New()

	handlers.ModifyPaths(app)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	port := os.Getenv("PORT")
	app.Listen(fmt.Sprintf(":%s", port))
}
