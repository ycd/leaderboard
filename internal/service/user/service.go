package user

import (
	"log"

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

func (s *UserService) SetProfile(user *models.User) (*models.User, error) {
	existingUser, err := s.userRepo.GetUserByUsername(user.Username)
	if err != nil {
		if errors.IsErrorType(err, errors.ErrUserNotFound) {
			log.Printf("user not found, proceeding to create new user")
		} else {
			return nil, err
		}
	}

	if existingUser != nil {
		if user.ID == "" {
			return nil, errors.NewError(errors.ErrUserAlreadyExists, "user already exists")
		}
		log.Printf("existing user: %v", existingUser)
		// Ensure username and country cannot be changed once set
		if user.Username != existingUser.Username {
			return nil, errors.NewError(errors.ErrInvalidRequest, "Username cannot be changed once set")
		}
		if user.Country != existingUser.Country {
			return nil, errors.NewError(errors.ErrInvalidRequest, "Country cannot be changed once set")
		}

		// Preserve existing level and coin if not provided
		if user.Level == 0 {
			user.Level = existingUser.Level
		}
		if user.Coin == 0 {
			user.Coin = existingUser.Coin
		}

		log.Printf("updating user profile: %v", user)
		// Proceed to update other mutable fields if any
		err := s.userRepo.UpdateUserProfile(user)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	// Handle new user creation
	if user.ID == "" {
		log.Printf("creating new user: %v", user)
		uid, err := s.userRepo.CreateUser(user)
		if err != nil {
			return nil, err
		}
		user.ID = uid
		log.Printf("new user created: %v", user)
		return user, nil
	}

	return nil, errors.NewError(errors.ErrInvalidRequest, "Invalid user profile data")
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
