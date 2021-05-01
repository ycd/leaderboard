package leaderboard

import (
	"context"

	"github.com/ycd/leaderboard/pkg/storage"
)

type Leaderboard struct {
	storage *storage.Storage
}

func NewLeaderboard(ctx context.Context) *Leaderboard {
	return &Leaderboard{
		storage: storage.NewStorage(ctx),
	}
}

func (l *Leaderboard) GetLeaderboard() (interface{}, error) {
	data, err := l.storage.GetLeaderboard()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (l *Leaderboard) GetLeaderboardWithCountry(country string) (interface{}, error) {
	data, err := l.storage.GetLeaderboardWithCountry(country)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (l *Leaderboard) ScoreSubmit(body *ScoreSubmit) (interface{}, error) {
	data, err := l.storage.ScoreSubmit(body.ScoreWorth, body.UserID, body.Timestamp)
	if err != nil {
		return nil, err
	}

	return data, nil
}
