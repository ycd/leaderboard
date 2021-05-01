package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ycd/leaderboard/pkg/leaderboard"
)

func HandleLeaderboard(c *fiber.Ctx) error {
	if err := leaderboard.NewLeaderboard(c.Context()).GetLeaderboard(); err != nil {
		return c.Status(500).JSON(Response{
			Err:     err.Error(),
			Success: false,
		})
	}

	return c.Status(200).JSON(Response{
		Err:     "",
		Success: true,
		Data:    "Test",
	})
}
