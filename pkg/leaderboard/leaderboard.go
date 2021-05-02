package leaderboard

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ycd/leaderboard/pkg/storage"
)

// Leaderboard is an abstraction between the Data Access Layer
// and Service layer
// Leaderboard implements the common functionalites
// to get the Leaderboard API service work.
type Leaderboard struct {
	storage *storage.Storage
}

// NewLeaderboard creates a new Leaderboard instance.
func NewLeaderboard(ctx context.Context) *Leaderboard {
	return &Leaderboard{
		storage: storage.NewStorage(ctx),
	}
}

// Health is used to make health checks over PostgreSQL connection.
func (l *Leaderboard) Health(ctx context.Context) error {
	return l.storage.Conn().Ping(ctx)
}

// GetLeaderboard returns the whole leaderboard.
func (l *Leaderboard) GetLeaderboard() (interface{}, error) {
	data, err := l.storage.GetLeaderboard()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetLeaderboardWithCountry returns the leaderboard with users frrom specific country.
func (l *Leaderboard) GetLeaderboardWithCountry(country string) (interface{}, error) {
	data, err := l.storage.GetLeaderboardWithCountry(country)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ScoreSubmit submit's the user score record to the database.
func (l *Leaderboard) ScoreSubmit(body *ScoreSubmit) (interface{}, error) {
	data, err := l.storage.ScoreSubmit(body.ScoreWorth, body.UserID, time.Now().Unix())
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UserCreate creates a new user and inserts it to the database.
func (l *Leaderboard) UserCreate(body *UserCreate) (interface{}, error) {
	data, err := l.storage.UserCreate(uuid.NewString(), body.UserName, body.Country)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetUser creates a new user and inserts it to the database.
func (l *Leaderboard) GetUser(guid string) (interface{}, error) {
	data, err := l.storage.GetUser(guid)
	if err != nil {
		return nil, err
	}

	return data, nil
}
