package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ycd/leaderboard/internal/errors"
	"github.com/ycd/leaderboard/internal/models"
	"github.com/ycd/leaderboard/internal/service/event"
)

type EventHandler struct {
	eventService *event.EventService
}

func NewEventHandler(eventService *event.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

func (h *EventHandler) JoinEvent(c *fiber.Ctx) error {
	var req models.JoinEventRequest
	if err := c.BodyParser(&req); err != nil {
		return handleError(c, errors.NewError(errors.ErrInvalidRequest, err.Error()))
	}

	err := h.eventService.JoinEvent(req.EventID, req.UserID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.JoinEventResponse{
		UserID:  req.UserID,
		EventID: req.EventID,
		Success: true,
	})
}

func (h *EventHandler) SetLeaderboardProgress(c *fiber.Ctx) error {
	var req models.SetLeaderboardProgressRequest
	if err := c.BodyParser(&req); err != nil {
		return handleError(c, errors.NewError(errors.ErrInvalidRequest, err.Error()))
	}

	err := h.eventService.SetLeaderboardProgress(req.EventID, req.UserID, req.Score)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SetLeaderboardProgressResponse{
		UserID:  req.UserID,
		EventID: req.EventID,
		Score:   req.Score,
		Success: true,
	})
}

func (h *EventHandler) GetLeaderboard(c *fiber.Ctx) error {
	var req models.GetLeaderboardRequest
	if c.Query("event_id") == "" {
		return handleError(c, errors.NewError(errors.ErrInvalidRequest, "event_id is required"))
	} else {
		req.EventID = c.Query("event_id")
	}

	if c.Query("limit") == "" {
		req.Limit = 10
	} else {
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			return handleError(c, errors.NewError(errors.ErrInvalidRequest, "limit is invalid"))
		}
		req.Limit = limit
	}

	entries, err := h.eventService.GetLeaderboard(req.EventID, req.Limit)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.LeaderboardResponse{
		EventID: req.EventID,
		Entries: entries,
	})
}

func (h *EventHandler) GetCountryLeaderboard(c *fiber.Ctx) error {
	var req models.GetCountryLeaderboardRequest
	if c.Query("event_id") == "" {
		return handleError(c, errors.NewError(errors.ErrInvalidRequest, "event_id is required"))
	} else {
		req.EventID = c.Query("event_id")
	}

	if c.Query("country") == "" {
		return handleError(c, errors.NewError(errors.ErrInvalidRequest, "country is required"))
	} else {
		req.Country = c.Query("country")
	}

	if c.Query("limit") == "" {
		req.Limit = 10
	} else {
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			return handleError(c, errors.NewError(errors.ErrInvalidRequest, "limit is invalid"))
		}
		req.Limit = limit
	}

	entries, err := h.eventService.GetCountryLeaderboard(req.EventID, req.Country, req.Limit)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.CountryLeaderboardResponse{
		EventID: req.EventID,
		Country: req.Country,
		Entries: entries,
	})
}

func (h *EventHandler) ClaimReward(c *fiber.Ctx) error {
	var req models.ClaimLeaderboardRequest
	if err := c.BodyParser(&req); err != nil {
		return handleError(c, errors.NewError(errors.ErrInvalidRequest, err.Error()))
	}

	reward, err := h.eventService.ClaimReward(req.EventID, req.UserID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.ClaimRewardResponse{
		UserID:     req.UserID,
		EventID:    req.EventID,
		RewardCoin: reward,
		Success:    true,
	})
}
