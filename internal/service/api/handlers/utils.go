package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ycd/leaderboard/internal/errors"
)

func handleError(c *fiber.Ctx, err error) error {
	log.Printf("Error: %v", err)
	apiError, ok := err.(errors.Error)
	if !ok {
		// If it's not our custom error type, treat it as an internal server error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	var statusCode int
	switch apiError.Type {
	case errors.ErrUserNotFound, errors.ErrEventNotFound, errors.ErrLeaderboardNotFound, errors.ErrCountryLeaderboardNotFound:
		statusCode = fiber.StatusNotFound
	case errors.ErrInvalidUsername, errors.ErrInvalidCountry, errors.ErrInvalidCoin, errors.ErrInvalidScore, errors.ErrInvalidEventDates, errors.ErrInvalidRequest:
		statusCode = fiber.StatusBadRequest
	case errors.ErrUnauthorized, errors.ErrUserBanned, errors.ErrInvalidLevel:
		statusCode = fiber.StatusUnauthorized
	case errors.ErrInsufficientLevel, errors.ErrUserAlreadyJoined, errors.ErrEventEnded, errors.ErrEventNotStarted, errors.ErrNotEligibleForReward, errors.ErrRewardAlreadyClaimed, errors.ErrUserAlreadyExists:
		statusCode = fiber.StatusForbidden
	default:
		statusCode = fiber.StatusInternalServerError
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"error":   apiError.Type,
		"message": apiError.Message,
	})
}
