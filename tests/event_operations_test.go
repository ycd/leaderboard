package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ycd/leaderboard/internal/models"
	"github.com/ycd/leaderboard/internal/repository"
	"github.com/ycd/leaderboard/internal/service/event"
	"github.com/ycd/leaderboard/internal/service/user"
	"github.com/ycd/leaderboard/internal/storage"
)

func TestEventOperations(t *testing.T) {
	sqliteDB, err := storage.NewSQLiteConnection()
	if err != nil {
		t.Fatalf("Failed to connect to SQLite: %v", err)
	}
	defer sqliteDB.Close()

	db, err := storage.NewClickHouseConnection()
	if err != nil {
		t.Fatalf("Failed to connect to ClickHouse: %v", err)
	}

	redisDB, err := storage.NewRedisConnection()
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisDB.Close()

	userRepo := repository.NewUserRepository(sqliteDB)
	eventRepo := repository.NewEventRepository(db, sqliteDB, redisDB) // Added sqliteDB as the second argument
	adminRepo := repository.NewAdminRepository(db, sqliteDB)
	userService := user.NewUserService(userRepo)
	eventService := event.NewEventService(eventRepo, userRepo)

	// Create test users
	user1 := createTestUser(t, userService, "user1", "US", 85)
	user2 := createTestUser(t, userService, "user2", "UK", 90)
	user3 := createTestUser(t, userService, "user3", "US", 75)

	// Create test event
	eventID := createTestEvent(t, adminRepo)

	// Test JoinEvent
	t.Run("JoinEvent", func(t *testing.T) {
		err := eventService.JoinEvent(eventID, user1.ID)
		assert.NoError(t, err)

		err = eventService.JoinEvent(eventID, user2.ID)
		assert.NoError(t, err)

		err = eventService.JoinEvent(eventID, user3.ID)
		assert.Error(t, err) // Should fail due to low level
	})

	// Test SetLeaderboardProgress
	t.Run("SetLeaderboardProgress", func(t *testing.T) {
		err := eventService.SetLeaderboardProgress(eventID, user1.ID, 100)
		assert.NoError(t, err)

		err = eventService.SetLeaderboardProgress(eventID, user2.ID, 200)
		assert.NoError(t, err)
	})

	// Test GetLeaderboard
	t.Run("GetLeaderboard", func(t *testing.T) {
		leaderboard, err := eventService.GetLeaderboard(eventID, 10)
		assert.NoError(t, err)
		assert.Len(t, leaderboard, 2)
		assert.Equal(t, user2.ID, leaderboard[0].UserID)
		assert.Equal(t, user1.ID, leaderboard[1].UserID)
	})

	// Test GetCountryLeaderboard
	t.Run("GetCountryLeaderboard", func(t *testing.T) {
		leaderboard, err := eventService.GetCountryLeaderboard(eventID, "US", 10)
		assert.NoError(t, err)
		assert.Len(t, leaderboard, 1)
		assert.Equal(t, user1.ID, leaderboard[0].UserID)
	})

	// Test ClaimReward
	t.Run("ClaimReward", func(t *testing.T) {
		reward, err := eventService.ClaimReward(eventID, user2.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1000, reward)

		reward, err = eventService.ClaimReward(eventID, user1.ID)
		assert.NoError(t, err)
		assert.Equal(t, 500, reward)

		_, err = eventService.ClaimReward(eventID, user3.ID)
		assert.Error(t, err) // Should fail as user3 is not in top 3
	})
}

func createTestUser(t *testing.T, userService *user.UserService, username, country string, level int) *models.User {
	user := &models.User{
		Username: username,
		Country:  country,
	}

	err := userService.SetProfile(user)
	assert.NoError(t, err)

	return user
}

func createTestEvent(t *testing.T, adminRepo *repository.AdminRepository) string {
	event := &models.Event{
		Name:      "Test Event",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
	}
	err := adminRepo.CreateEvent(event)
	assert.NoError(t, err)
	return event.ID
}
