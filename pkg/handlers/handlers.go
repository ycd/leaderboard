package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ycd/leaderboard/pkg/leaderboard"
)

func HandleLeaderboard(c *fiber.Ctx) error {
	data, err := leaderboard.NewLeaderboard(c.Context()).GetLeaderboard()
	if err != nil {
		return c.Status(500).JSON(Response{
			Err:     err.Error(),
			Success: false,
			Data:    nil,
		})
	}

	return c.Status(200).JSON(Response{
		Err:     "",
		Success: true,
		Data:    data,
	})
}

func HandleLeaderboardWithcountry(c *fiber.Ctx) error {
	country := c.Query("country")
	data, err := leaderboard.NewLeaderboard(c.Context()).GetLeaderboardWithCountry(country)
	if err != nil {
		return c.Status(500).JSON(Response{
			Err:     err.Error(),
			Success: false,
			Data:    nil,
		})
	}

	return c.Status(200).JSON(Response{
		Err:     "",
		Success: true,
		Data:    data,
	})
}

func HandleScoreSubmit(c *fiber.Ctx) error {
	b := new(leaderboard.ScoreSubmit)
	if err := c.BodyParser(b); err != nil {
		return err
	}

	data, err := leaderboard.NewLeaderboard(c.Context()).ScoreSubmit(b)
	if err != nil {
		return c.Status(500).JSON(Response{
			Err:     err.Error(),
			Success: false,
			Data:    nil,
		})
	}

	return c.Status(200).JSON(Response{
		Err:     "",
		Success: true,
		Data:    data,
	})
}
