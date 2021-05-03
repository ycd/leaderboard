package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/ycd/leaderboard/pkg/leaderboard"
)

// Handler is an abstraction for
// accessing the leaderboard
// and sharing the underlying connection pool
// in every handler and requests.
type Handler struct {
	leaderboard *leaderboard.Leaderboard
}

// NewHandlers returns a new Handler instance.
func NewHandler(ctx context.Context) *Handler {
	return &Handler{
		leaderboard: leaderboard.NewLeaderboard(ctx),
	}
}

// HandleHealthCheck handler pings the database and conforms the liveness probe.
func (h *Handler) HandleHealthCheck(c *fiber.Ctx) error {
	ctx := context.Background()
	if err := h.leaderboard.Health(ctx); err != nil {
		return c.Status(503).Send([]byte(err.Error()))
	}

	return c.Status(200).Send([]byte("ok"))
}

// HandleLeaderboard retrieves the global leaderboard.
func (h *Handler) HandleLeaderboard(c *fiber.Ctx) error {
	data, err := h.leaderboard.GetLeaderboard()
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

// HandleLeaderboardWithcountry retrieves the leaderboard of the given country.
func (h *Handler) HandleLeaderboardWithcountry(c *fiber.Ctx) error {
	country := c.Params("country")

	data, err := h.leaderboard.GetLeaderboardWithCountry(country)
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

// HandleScoreSubmit submits a new score record to the database.
func (h *Handler) HandleScoreSubmit(c *fiber.Ctx) error {
	b := new(leaderboard.ScoreSubmit)
	if err := c.BodyParser(b); err != nil {
		return err
	}

	data, err := h.leaderboard.ScoreSubmit(b)
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
func (h *Handler) HandleUserCreate(c *fiber.Ctx) error {
	b := new(leaderboard.UserCreate)
	if err := c.BodyParser(b); err != nil {
		return err
	}

	data, err := h.leaderboard.UserCreate(b)
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

// HandleGetUser retrieves the information of the user with given ID.
func (h *Handler) HandleGetUser(c *fiber.Ctx) error {
	guid := c.Params("guid")

	data, err := h.leaderboard.GetUser(guid)
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
