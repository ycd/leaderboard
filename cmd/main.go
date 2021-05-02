package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/ycd/leaderboard/pkg/handlers"
)

func main() {
	app := fiber.New()
	handlers.ModifyPaths(app)

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
