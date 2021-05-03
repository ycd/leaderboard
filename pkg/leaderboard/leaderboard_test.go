package leaderboard

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/ycd/leaderboard/pkg/storage"
)

var ids = []string{}

func init() {
	ctx := context.Background()
	cleanTables()
	NewLeaderboard(ctx).storage.CreateTables(ctx)
}

func TestNewLeaderboard(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name string
	}{
		{
			name: "new-leaderboard",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewLeaderboard(ctx).Health(ctx); err != nil {
				t.Errorf("Health check failed: %v", err)
			}
		})
	}
}

func TestLeaderboardHealth(t *testing.T) {
	ctx := context.Background()
	leaderboard := NewLeaderboard(ctx)
	tests := []struct {
		name string
	}{
		{
			name: "health-check",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := leaderboard.Health(ctx); err != nil {
				t.Errorf("Health check failed: %v", err)
			}
		})
	}
}

func TestLeaderboardUserCreate(t *testing.T) {
	ctx := context.Background()
	leaderboard := NewLeaderboard(ctx)
	tests := []struct {
		name string
		body UserCreate
	}{
		{
			name: "test-user-create",
			body: UserCreate{
				UserName: "test_araba",
				Country:  "br",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := leaderboard.UserCreate(&tt.body)
			if err != nil {
				t.Fatalf("user creation failed: %v", err)
			}

			user := got.(storage.UserInfo)
			ids = append(ids, user.UserID)

			if user.DisplayName != tt.body.UserName {
				t.Fatalf("expected name: %v, got: %v", tt.body.UserName, user.DisplayName)
			}

			if user.Points != 0 {
				t.Fatalf("newly created user msut have point of: 0, got: %d", user.Points)
			}
		})
	}
}

func TestLeaderboardScoreSubmit(t *testing.T) {
	ctx := context.Background()
	leaderboard := NewLeaderboard(ctx)
	tests := []struct {
		name string
		body ScoreSubmit
	}{
		{
			name: "test-submit-score",
			body: ScoreSubmit{
				ScoreWorth: 22,
				UserID:     ids[0],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := leaderboard.ScoreSubmit(&tt.body)
			if err != nil {
				t.Fatalf("score submission failed: %v", err)
			}

			score := got.(storage.Score)

			if score.UserID != tt.body.UserID {
				t.Fatalf("expected user with id: %v, got: %v", tt.body.UserID, score.UserID)
			}
		})
	}
}

func TestLeaderboardGetUser(t *testing.T) {
	ctx := context.Background()
	leaderboard := NewLeaderboard(ctx)
	tests := []struct {
		name   string
		userID string
	}{
		{
			name:   "test-get-user",
			userID: ids[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := leaderboard.GetUser(tt.userID)
			if err != nil {
				t.Fatalf("failed to retrieve user data: %v", err)
			}

			user := got.(storage.UserInfo)

			if user.Points != 22 {
				t.Fatalf("expected user to be have %d points, got: %v", 22, user.Points)
			}
		})
	}
}

func TestLeaderboardGetLeaderboard(t *testing.T) {
	ctx := context.Background()
	leaderboard := NewLeaderboard(ctx)
	tests := []struct {
		name    string
		wantErr bool
		want    []storage.LeaderboardResult
	}{
		{
			name:    "test-get-leaderboard",
			wantErr: false,
			want: []storage.LeaderboardResult{
				{
					Rank:        1,
					Points:      22,
					DisplayName: "test_araba",
					Country:     "br",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := leaderboard.GetLeaderboard()
			if err != nil {
				t.Fatalf("failed to retrieve leaderboard: %v", err)
			}
			var lb []storage.LeaderboardResult

			js, err := json.Marshal(got)
			if err != nil {
				log.Fatalf("got error while marshalling: %v", err)
			}
			if err := json.Unmarshal(js, &lb); err != nil {
				log.Fatalf("got error while unmarshalling into struct: %v", err)
			}

			if lb[0] != tt.want[0] {
				log.Fatalf("expected: %v, got: %v", tt.want, lb)
			}
		})
	}
}

func TestLeaderboardGetLeaderboardWithCountry(t *testing.T) {
	ctx := context.Background()
	leaderboard := NewLeaderboard(ctx)
	defer t.Cleanup(cleanTables)

	tests := []struct {
		name    string
		wantErr bool
		want    []storage.LeaderboardResult
		country string
	}{
		{
			name:    "test-get-leaderboard",
			wantErr: false,
			want: []storage.LeaderboardResult{
				{
					Rank:        1,
					Points:      22,
					DisplayName: "test_araba",
					Country:     "br",
				},
			},
			country: "br",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := leaderboard.GetLeaderboardWithCountry(tt.country)
			if err != nil {
				t.Fatalf("failed to retrieve leaderboard: %v", err)
			}
			var lb []storage.LeaderboardResult

			js, err := json.Marshal(got)
			if err != nil {
				log.Fatalf("got error while marshalling: %v", err)
			}
			if err := json.Unmarshal(js, &lb); err != nil {
				log.Fatalf("got error while unmarshalling into struct: %v", err)
			}

			if lb[0] != tt.want[0] {
				log.Fatalf("expected: %v, got: %v", tt.want, lb)
			}
		})
	}
}

func cleanTables() {
	ctx := context.Background()
	lb := NewLeaderboard(ctx)
	tableNames := []string{
		"scores",
		"users",
	}
	viewNames := []string{
		"leaderboard",
		"UsersWithScores",
	}

	for _, tName := range viewNames {
		if _, err := lb.storage.Conn().Exec(ctx, fmt.Sprintf("DROP VIEW %s", tName)); err != nil {
			log.Printf("failed to delete view table :%v", err)
		}
		log.Println("leaderboard:", tName)
	}

	for _, tName := range tableNames {
		if _, err := lb.storage.Conn().Exec(ctx, fmt.Sprintf("DROP TABLE %s", tName)); err != nil {
			log.Printf("failed to drop table :%v", err)
		}
		log.Println("leaderboard", tName)
	}

	lb.storage.Cache().Del(ctx, "leaderboard")
}
