package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/ycd/leaderboard/pkg/handlers"
)

func main() {
	app := fiber.New()

	h := handlers.NewHandler(context.Background())
	h.ModifyPaths(app)

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
