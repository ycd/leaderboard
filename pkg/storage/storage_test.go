package storage

import (
	"context"
	"fmt"
	"log"
	"testing"
)

// Helper function for tests.
func MockStorage(ctx context.Context) *Storage {
	return NewStorage(ctx)
}

func TestNewStorage(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test-connection",
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewStorage(tt.args.ctx)
			if err := got.connection.Ping(tt.args.ctx); err != nil {
				t.Fatalf("NewStorage failed: %v", err)
			}
		})
	}
}

func TestConn(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		ctx     context.Context
	}{
		{
			name:    "test-conn",
			wantErr: false,
			ctx:     context.Background(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MockStorage(tt.ctx).Conn().Ping(tt.ctx); err != nil {
				t.Fatalf("Accessing underlying connection failed: %v", err)
			}
		})
	}
}

func TestCreateTables(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name    string
		storage *Storage
	}{
		{
			name:    "create-tables",
			storage: MockStorage(ctx),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storage.CreateTables(ctx)

			tableNames := []string{"scores", "users", "UsersWithScores", "leaderboard"}

			for _, tName := range tableNames {
				if _, err := tt.storage.Conn().Exec(ctx, fmt.Sprintf("SELECT 1 FROM %s ", tName)); err != nil {
					t.Fatalf("Table creation failed for: %v", err)
				}
			}

			cleanTables()
		})
	}
}

func cleanTables() {
	ctx := context.Background()
	storage := MockStorage(ctx)
	tableNames := []string{
		"scores",
		"users",
	}
	viewNames := []string{
		"leaderboard",
		"UsersWithScores",
	}

	for _, tName := range viewNames {
		if _, err := storage.Conn().Exec(ctx, fmt.Sprintf("DROP	 VIEW %s", tName)); err != nil {
			log.Printf("failed to delete view table :%v", err)
		}
	}

	for _, tName := range tableNames {
		if _, err := storage.Conn().Exec(ctx, fmt.Sprintf("DROP TABLE %s", tName)); err != nil {
			log.Printf("failed to drop table :%v", err)
		}
	}

	storage.Cache().Del(ctx, "leaderboard")
}
