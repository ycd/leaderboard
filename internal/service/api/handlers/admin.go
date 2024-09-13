package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ycd/leaderboard/internal/errors"
	"github.com/ycd/leaderboard/internal/models"
	"github.com/ycd/leaderboard/internal/service/admin"
)

type AdminHandler struct {
	adminService *admin.AdminService
}

func NewAdminHandler(adminService *admin.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) CreateEvent(c *fiber.Ctx) error {
	var req models.CreateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return handleError(c, errors.NewError(errors.ErrInvalidRequest, err.Error()))
	}

	event := &models.Event{
		Name:      req.Name,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	event, err := h.adminService.CreateEvent(event)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(models.CreateEventResponse{Event: *event})
}

func (h *AdminHandler) ListEvents(c *fiber.Ctx) error {
	events, err := h.adminService.ListEvents()
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.ListEventsResponse{Events: events})
}

func (h *AdminHandler) DeleteEvent(c *fiber.Ctx) error {
	var req models.DeleteEventRequest
	if err := c.BodyParser(&req); err != nil {
		return handleError(c, errors.NewError(errors.ErrInvalidRequest, err.Error()))
	}

	err := h.adminService.DeleteEvent(req.EventID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.DeleteEventResponse{Success: true})
}

func (h *AdminHandler) GetUserDetails(c *fiber.Ctx) error {
	var req models.GetUserDetailsRequest
	if c.Query("user_id") == "" {
		return handleError(c, errors.NewError(errors.ErrInvalidRequest, "Invalid user ID"))
	} else {
		req.UserID = c.Query("user_id")
	}

	user, err := h.adminService.GetUserDetails(req.UserID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.GetUserDetailsResponse{User: *user})
}

func (h *AdminHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.adminService.ListUsers()
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.ListUsersResponse{Users: users})
}

func (h *AdminHandler) BanUser(c *fiber.Ctx) error {
	var req models.BanUserRequest
	if err := c.BodyParser(&req); err != nil {
		return handleError(c, errors.NewError(errors.ErrInvalidRequest, err.Error()))
	}

	err := h.adminService.BanUser(req.UserID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.BanUserResponse{Success: true})
}

func (h *AdminHandler) UnbanUser(c *fiber.Ctx) error {
	var req models.UnbanUserRequest
	if err := c.BodyParser(&req); err != nil {
		return handleError(c, errors.NewError(errors.ErrInvalidRequest, err.Error()))
	}

	err := h.adminService.UnbanUser(req.UserID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.UnbanUserResponse{Success: true})
}

func (h *AdminHandler) ResetLeaderboard(c *fiber.Ctx) error {
	var req models.ResetLeaderboardRequest
	if err := c.BodyParser(&req); err != nil {
		return handleError(c, errors.NewError(errors.ErrInvalidRequest, err.Error()))
	}

	err := h.adminService.ResetLeaderboard(req.EventID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.ResetLeaderboardResponse{Success: true})
}
