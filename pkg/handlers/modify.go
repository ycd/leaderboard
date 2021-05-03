package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// ModifyPaths registers the handlers to the main app.
func (h *Handler) ModifyPaths(f *fiber.App) {
	f.Get("/health", h.HandleHealthCheck)

	f.Get("/leaderboard", h.HandleLeaderboard)
	f.Get("/leaderboard/:country", h.HandleLeaderboardWithcountry)
	f.Get("/user/profile/:guid", h.HandleGetUser)
	f.Post("/user/create", h.HandleUserCreate)
	f.Post("/score/submit", h.HandleScoreSubmit)
}
