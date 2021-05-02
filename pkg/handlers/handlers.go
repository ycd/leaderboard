package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/ycd/leaderboard/pkg/leaderboard"
)

// HandleHealthCheck handler pings the database and conforms the liveness probe.
func HandleHealthCheck(c *fiber.Ctx) error {
	ctx := context.Background()
	if err := leaderboard.NewLeaderboard(ctx).Health(ctx); err != nil {
		return c.Status(503).Send([]byte(err.Error()))
	}

	return c.Status(200).Send([]byte("ok"))
}

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
		Success: true,
		Data:    data,
	})
}

func HandleLeaderboardWithcountry(c *fiber.Ctx) error {
	country := c.Params("country")

	data, err := leaderboard.NewLeaderboard(c.Context()).GetLeaderboardWithCountry(country)
	if err != nil {
		return c.Status(500).JSON(Response{
			Err:     err.Error(),
			Success: false,
			Data:    nil,
		})
	}

	return c.Status(200).JSON(Response{
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
		Success: true,
		Data:    data,
	})
}

// HandleUserCreate handles the logic for user creation
//
// On success returns HTTP 201 - Created
// On failure - or payload with duplicate name returns HTTP 409 - Conflict
func HandleUserCreate(c *fiber.Ctx) error {
	b := new(leaderboard.UserCreate)
	if err := c.BodyParser(b); err != nil {
		return err
	}

	data, err := leaderboard.NewLeaderboard(c.Context()).UserCreate(b)
	if err != nil {
		return c.Status(409).JSON(Response{
			Err:     err.Error(),
			Success: false,
			Data:    nil,
		})
	}

	return c.Status(201).JSON(Response{
		Success: true,
		Data:    data,
	})
}

func HandleGetUser(c *fiber.Ctx) error {
	guid := c.Params("guid")

	data, err := leaderboard.NewLeaderboard(c.Context()).GetUser(guid)
	if err != nil {
		return c.Status(500).JSON(Response{
			Err:     err.Error(),
			Success: false,
			Data:    nil,
		})
	}

	return c.Status(200).JSON(Response{
		Success: true,
		Data:    data,
	})
}
