package event

import (
	"context"
	"time"

	"github.com/ycd/leaderboard/internal/errors"
	"github.com/ycd/leaderboard/internal/models"
	"github.com/ycd/leaderboard/internal/repository"
)

type EventService struct {
	eventRepo *repository.EventRepository
	userRepo  *repository.UserRepository
}

func NewEventService(eventRepo *repository.EventRepository, userRepo *repository.UserRepository) *EventService {
	return &EventService{
		eventRepo: eventRepo,
		userRepo:  userRepo,
	}
}

func (s *EventService) JoinEvent(eventID, userID string) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.Level < 80 {
		return errors.NewError(errors.ErrInsufficientLevel, "User level must be at least 80 to join an event")
	}

	return s.eventRepo.AddParticipant(eventID, userID)
}

func (s *EventService) SetLeaderboardProgress(eventID, userID string, score int) error {
	return s.eventRepo.AddParticipantScore(eventID, userID, score)
}

func (s *EventService) GetLeaderboard(eventID string, limit int) ([]models.LeaderboardEntryWithRank, error) {
	return s.eventRepo.GetLeaderboard(eventID, limit)
}

func (s *EventService) GetCountryLeaderboard(eventID string, country string, limit int) ([]models.CountryLeaderboardEntryWithRank, error) {
	return s.eventRepo.GetCountryLeaderboard(eventID, country, limit)
}

func (s *EventService) ClaimReward(eventID, userID string) (int, error) {
	rank, err := s.eventRepo.GetParticipantRank(eventID, userID)
	if err != nil {
		return 0, err
	}

	var reward int
	switch rank {
	case 1:
		reward = 1000
	case 2:
		reward = 500
	case 3:
		reward = 250
	default:
		return 0, errors.NewError(errors.ErrNotEligibleForReward, "No reward is available")
	}

	event, err := s.eventRepo.GetEventByID(context.Background(), eventID)
	if err != nil {
		return 0, err
	}

	if event.EndTime.After(time.Now()) {
		return 0, errors.NewError(errors.ErrInvalidRequest, "Event is not yet ended")
	}

	err = s.eventRepo.ClaimReward(eventID, userID, reward)
	if err != nil {
		return 0, err
	}

	err = s.userRepo.AddUserCoin(userID, reward)
	if err != nil {
		return 0, err
	}

	return reward, nil
}
