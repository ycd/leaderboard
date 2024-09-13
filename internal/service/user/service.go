package user

import (
	"github.com/ycd/leaderboard/internal/errors"
	"github.com/ycd/leaderboard/internal/models"
	"github.com/ycd/leaderboard/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) SetProfile(user *models.User) error {
	existingUser, err := s.userRepo.GetUserByUsername(user.Username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.NewError(errors.ErrUserAlreadyExists, "Username already exists")
	}

	if user.ID == "" {
		usr, err := s.userRepo.CreateUser(user)
		if err != nil {
			return err
		}

		user.ID = usr
		return nil
	}

	return nil
}

func (s *UserService) SetLevel(userID string, level int) error {
	if level < 0 {
		return errors.NewError(errors.ErrInvalidLevel, "Level must be non-negative")
	}
	return s.userRepo.SetUserLevel(userID, level)
}

func (s *UserService) SetCoin(userID string, coin int) error {
	return s.userRepo.SetUserCoin(userID, coin)
}
