package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ycd/leaderboard/internal/errors"
	"github.com/ycd/leaderboard/internal/models"
	"github.com/ycd/leaderboard/internal/service/user"
)

type UserHandler struct {
	userService *user.UserService
}

func NewUserHandler(userService *user.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) SetProfile(c *fiber.Ctx) error {
	log.Println("SetProfile called")
	var req models.SetProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return handleError(c, errors.NewError(errors.ErrInvalidRequest, err.Error()))
	}

	user := &models.User{
		Username: req.Username,
		Country:  req.Country,
		Level:    int32(req.Level),
		Coin:     int32(req.Coin),
	}

	usr, err := h.userService.SetProfile(user)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.UserProfileResponse{
		ID:       usr.ID,
		Username: usr.Username,
		Country:  usr.Country,
		Level:    int(usr.Level),
		Coin:     int(usr.Coin),
	})
}

func (h *UserHandler) SetLevel(c *fiber.Ctx) error {
	log.Println("SetLevel called")
	var req models.SetLevelRequest
	if err := c.BodyParser(&req); err != nil {
		return handleError(c, errors.NewError(errors.ErrInvalidRequest, err.Error()))
	}

	err := h.userService.SetLevel(req.UserID, req.Level)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SetLevelResponse{
		UserID: req.UserID,
		Level:  req.Level,
	})
}

func (h *UserHandler) SetCoin(c *fiber.Ctx) error {
	log.Println("SetCoin called")
	var req models.SetCoinRequest
	if err := c.BodyParser(&req); err != nil {
		return handleError(c, errors.NewError(errors.ErrInvalidRequest, err.Error()))
	}

	err := h.userService.SetCoin(req.UserID, req.Coin)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SetCoinResponse{
		UserID: req.UserID,
		Coin:   req.Coin,
	})
}
