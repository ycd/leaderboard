package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func ModifyPaths(f *fiber.App) {
	f.Get("/leaderboard", HandleLeaderboard)
	// f.Get("/leaderboard/:country")
	// f.Get("/user/profile/:guid")
	// f.Post("/user/create")
	// f.Post("/score/submit")
}
