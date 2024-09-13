package admin

import (
	"github.com/ycd/leaderboard/internal/models"
	"github.com/ycd/leaderboard/internal/repository"
)

type AdminService struct {
	userRepo  *repository.UserRepository
	eventRepo *repository.EventRepository
	adminRepo *repository.AdminRepository
}

func NewAdminService(userRepo *repository.UserRepository, eventRepo *repository.EventRepository, adminRepo *repository.AdminRepository) *AdminService {
	return &AdminService{
		userRepo:  userRepo,
		eventRepo: eventRepo,
		adminRepo: adminRepo,
	}
}

func (s *AdminService) CreateEvent(event *models.Event) (*models.Event, error) {
	err := s.adminRepo.CreateEvent(event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *AdminService) ListEvents() ([]models.Event, error) {
	return s.adminRepo.ListEvents()
}

func (s *AdminService) DeleteEvent(eventID string) error {
	return s.adminRepo.DeleteEvent(eventID)
}

func (s *AdminService) GetUserDetails(userID string) (*models.User, error) {
	return s.userRepo.GetUserByID(userID)
}

func (s *AdminService) ListUsers() ([]models.User, error) {
	return s.adminRepo.ListUsers()
}

func (s *AdminService) BanUser(userID string) error {
	return s.adminRepo.BanUser(userID)
}

func (s *AdminService) UnbanUser(userID string) error {
	return s.adminRepo.UnbanUser(userID)
}

func (s *AdminService) ResetLeaderboard(eventID string) error {
	return s.adminRepo.ResetLeaderboard(eventID)
}
