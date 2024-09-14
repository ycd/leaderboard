package tests

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/ycd/leaderboard/internal/models"
	"github.com/ycd/leaderboard/internal/repository"
	"github.com/ycd/leaderboard/internal/service/event"
	"github.com/ycd/leaderboard/internal/service/user"
	"github.com/ycd/leaderboard/internal/storage"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("No .env file found. Proceeding with environment variables.")
	}
}

func TestEventOperations(t *testing.T) {
	sqliteDB, err := storage.NewSQLiteConnection()
	if err != nil {
		t.Fatalf("Failed to connect to SQLite: %v", err)
	}
	defer sqliteDB.Close()

	chDB, err := storage.NewClickHouseConnection()
	if err != nil {
		t.Fatalf("Failed to connect to ClickHouse: %v", err)
	}

	redisDB, err := storage.NewRedisConnection()
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisDB.Close()

	userRepo := repository.NewUserRepository(sqliteDB)
	eventRepo := repository.NewEventRepository(chDB, sqliteDB, redisDB)
	adminRepo := repository.NewAdminRepository(chDB, sqliteDB)
	userService := user.NewUserService(userRepo)
	eventService := event.NewEventService(eventRepo, userRepo)

	// Create test users

	user1 := createTestUser(t, userService, fmt.Sprintf("user-%s", uuid.New().String()), "US", 85)
	user2 := createTestUser(t, userService, fmt.Sprintf("user-%s", uuid.New().String()), "UK", 90)
	user3 := createTestUser(t, userService, fmt.Sprintf("user-%s", uuid.New().String()), "US", 115)

	// Create test event
	eventID := createTestEvent(t, adminRepo)
	log.Printf("Event ID: %s", eventID)
	log.Printf("User 1 ID: %s", user1.ID)
	log.Printf("User 2 ID: %s", user2.ID)
	log.Printf("User 3 ID: %s", user3.ID)

	t.Run("SetUserLevels", func(t *testing.T) {
		err := userService.SetLevel(user1.ID, 85)
		assert.NoError(t, err)

		err = userService.SetLevel(user2.ID, 90)
		assert.NoError(t, err)

		err = userService.SetLevel(user3.ID, 155)
		assert.NoError(t, err)
	})

	// Test JoinEvent
	t.Run("JoinEvent", func(t *testing.T) {
		err := eventService.JoinEvent(eventID, user1.ID)
		assert.NoError(t, err)

		err = eventService.JoinEvent(eventID, user2.ID)
		assert.NoError(t, err)

		err = eventService.JoinEvent(eventID, user3.ID)
		assert.NoError(t, err) // Should fail due to low level
	})

	// Test SetLeaderboardProgress
	t.Run("SetLeaderboardProgress", func(t *testing.T) {
		err := eventService.SetLeaderboardProgress(eventID, user1.ID, 100)
		assert.NoError(t, err)

		err = eventService.SetLeaderboardProgress(eventID, user2.ID, 200)
		assert.NoError(t, err)

		err = eventService.SetLeaderboardProgress(eventID, user3.ID, 300)
		assert.NoError(t, err)

		err = eventService.SetLeaderboardProgress(eventID, user3.ID, 300)
		assert.NoError(t, err)

		leaderboard, err := eventService.GetLeaderboard(eventID, 10)
		assert.NoError(t, err)
		assert.Len(t, leaderboard, 3)
		assert.Equal(t, user3.ID, leaderboard[0].UserID)
		assert.Equal(t, user2.ID, leaderboard[1].UserID)
		assert.Equal(t, user1.ID, leaderboard[2].UserID)
	})

	// Test GetLeaderboard
	t.Run("GetLeaderboard", func(t *testing.T) {
		leaderboard, err := eventService.GetLeaderboard(eventID, 10)
		assert.NoError(t, err)
		assert.Len(t, leaderboard, 3)
		assert.Equal(t, user3.ID, leaderboard[0].UserID)
		assert.Equal(t, user2.ID, leaderboard[1].UserID)
		assert.Equal(t, user1.ID, leaderboard[2].UserID)
	})

	// Test GetCountryLeaderboard
	t.Run("GetCountryLeaderboard", func(t *testing.T) {
		leaderboard, err := eventService.GetCountryLeaderboard(eventID, "US", 10)
		assert.NoError(t, err)
		assert.Len(t, leaderboard, 2)
		assert.Equal(t, user3.ID, leaderboard[0].UserID)
		assert.Equal(t, user1.ID, leaderboard[1].UserID)
	})

	// Test ClaimReward
	t.Run("ClaimReward", func(t *testing.T) {
		reward, _ := eventService.ClaimReward(eventID, user2.ID)

		assert.Equal(t, 0, reward)

		reward, _ = eventService.ClaimReward(eventID, user1.ID)

		assert.Equal(t, 0, reward)
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
