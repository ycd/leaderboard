package tests

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ycd/leaderboard/internal/models"
	"github.com/ycd/leaderboard/internal/repository"
	"github.com/ycd/leaderboard/internal/service/user"
	"github.com/ycd/leaderboard/internal/storage"
)

var db *sql.DB

func init() {
	chDB, err := storage.NewClickHouseConnection()
	if err != nil {
		panic(err)
	}
	db = chDB
}

func TestUserOperations(t *testing.T) {
	sqliteDB, err := storage.NewSQLiteConnection()
	if err != nil {
		t.Fatalf("Failed to connect to SQLite: %v", err)
	}
	defer sqliteDB.Close()

	userRepo := repository.NewUserRepository(sqliteDB)
	userService := user.NewUserService(userRepo)

	testUser := createTestUser(t, userService, "testuser-123", "US", 0)
	t.Run("SetProfile", func(t *testing.T) {
		user := &models.User{
			Username: "testuser-234",
			Country:  "US",
		}
		err := userService.SetProfile(user)
		assert.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.Equal(t, 0, user.Coin)
	})

	// Test SetLevel
	t.Run("SetLevel", func(t *testing.T) {
		err := userService.SetLevel(testUser.ID, 85)
		assert.NoError(t, err)

		updatedUser, err := userRepo.GetUserByID(testUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, 85, updatedUser.Level)
	})

	// Test SetCoin
	t.Run("SetCoin", func(t *testing.T) {
		err := userService.SetCoin(testUser.ID, 100)
		assert.NoError(t, err)

		updatedUser, err := userRepo.GetUserByID(testUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, 100, updatedUser.Coin)
	})
}
