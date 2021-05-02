package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func ModifyPaths(f *fiber.App) {
	f.Get("/leaderboard", HandleLeaderboard)
	f.Get("/leaderboard/:country", HandleLeaderboardWithcountry)
	// f.Get("/user/profile/:guid")
	f.Post("/user/create", HandleUserCreate)
	f.Post("/score/submit", HandleScoreSubmit)
}
