package storage

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
)

var ids = []string{}

func init() {
	ctx := context.Background()
	cleanTables()
	MockStorage(ctx).CreateTables(ctx)
}

func TestUserCreate(t *testing.T) {
	ctx := context.Background()
	storage := MockStorage(ctx)
	type args struct {
		userID   string
		userName string
		country  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "user-create-1",
			args: args{
				userID:   uuid.NewString(),
				userName: "ycd_1",
				country:  "tr",
			},
			wantErr: false,
		},
		{
			name: "user-create-2",
			args: args{
				userID:   uuid.NewString(),
				userName: "ycd_2",
				country:  "fr",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := storage.UserCreate(tt.args.userID, tt.args.userName, tt.args.country)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.UserCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			user := resp.(UserInfo)
			ids = append(ids, user.UserID)

			if user.Points != 0 {
				t.Fatalf("newly created user point must be: 0, got %d", user.Points)
			}

			if user.DisplayName != tt.args.userName {
				t.Fatalf("expected display_name: %s, got %s", tt.args.userName, user.DisplayName)
			}
		})
	}
}

func TestScoreSubmit(t *testing.T) {
	ctx := context.Background()
	storage := MockStorage(ctx)
	type args struct {
		ScoreWorth float64
		UserID     string
		Timestamp  int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "submit-score-user-1",
			args: args{
				ScoreWorth: 50,
				UserID:     ids[0],
				Timestamp:  time.Now().Unix(),
			},
			wantErr: false,
		},
		{
			name: "submit-score-user-2",
			args: args{
				ScoreWorth: 25,
				UserID:     ids[1],
				Timestamp:  time.Now().Unix(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := storage.ScoreSubmit(tt.args.ScoreWorth, tt.args.UserID, tt.args.Timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.UserCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			user := resp.(Score)

			if user.UserID != tt.args.UserID {
				t.Fatalf("expected user id %s, got: %s", tt.args.UserID, user.UserID)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	storage := MockStorage(ctx)

	tests := []struct {
		name          string
		id            string
		wantErr       bool
		expectedScore float64
	}{
		{
			name:          "get-user-1",
			id:            ids[0],
			wantErr:       false,
			expectedScore: 50,
		},
		{
			name:          "get-user-2",
			id:            ids[1],
			wantErr:       false,
			expectedScore: 25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := storage.GetUser(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.UserCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			user := resp.(UserInfo)

			if float64(user.Points) != tt.expectedScore {
				t.Fatalf("expected score: %f, got :%d", tt.expectedScore, user.Points)
			}
		})
	}
}

func TestGetLeaderboard(t *testing.T) {
	ctx := context.Background()
	storage := MockStorage(ctx)

	tests := []struct {
		name    string
		want    []LeaderboardResult
		wantErr bool
	}{
		{
			name: "test-leaderboard",
			want: []LeaderboardResult{
				{
					Rank:        1,
					Points:      50,
					DisplayName: "ycd_1",
					Country:     "tr",
				},
				{
					Rank:        2,
					Points:      25,
					DisplayName: "ycd_2",
					Country:     "fr",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leaderboard, err := storage.GetLeaderboard()
			if err != nil {
				t.Fatalf("got error: %v", err)
			}

			var lb []LeaderboardResult

			_, err = storage.CacheGet(ctx, "leaderboard")
			if err != nil {
				t.Fatalf("got error while accessing cache: %v", err)
			}

			js, err := json.Marshal(leaderboard)
			if err != nil {
				log.Fatalf("got error while marshalling: %v", err)
			}
			if err := json.Unmarshal(js, &lb); err != nil {
				log.Fatalf("got error while unmarshalling into struct: %v", err)
			}

			if lb[0] != tt.want[0] {
				log.Fatalf("expected: %v, got: %v", tt.want, lb)
			}
			if lb[1] != tt.want[1] {
				log.Fatalf("expected: %v, got: %v", tt.want, lb)
			}
		})
	}
}

func TestGetLeaderboardWithCountry(t *testing.T) {
	ctx := context.Background()
	storage := MockStorage(ctx)
	defer t.Cleanup(cleanTables)

	tests := []struct {
		name    string
		country string
		want    []LeaderboardResult
		wantErr bool
	}{
		{
			name:    "test-leaderboard",
			country: "fr",
			want: []LeaderboardResult{
				{
					Rank:        1,
					Points:      25,
					DisplayName: "ycd_2",
					Country:     "fr",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leaderboard, err := storage.GetLeaderboardWithCountry(tt.country)
			if err != nil {
				t.Fatalf("got error: %v", err)
			}

			var lb []LeaderboardResult

			_, err = storage.CacheGet(ctx, "leaderboard")
			if err != nil {
				t.Fatalf("got error while accessing cache: %v", err)
			}

			js, err := json.Marshal(leaderboard)
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
