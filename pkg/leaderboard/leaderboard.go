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

func (l *Leaderboard) GetLeaderboard() error {
	if err := l.storage.GetLeaderboard(); err != nil {
		return err
	}

	return nil
}
