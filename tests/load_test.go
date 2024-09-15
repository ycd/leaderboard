package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ycd/leaderboard/internal/models"
)

const (
	baseURL     = "http://leaderboard.yagizdegirmenci.com:8080"
	userCount   = 1000
	updateCount = 20000
)

func TestLoad(t *testing.T) {
	// Create users concurrently
	users := createUsers(t)

	// Set user levels concurrently
	setUserLevels(t, users)

	// Create an event
	event := createEvent(t)

	// Join users to the event concurrently
	joinUsersToEvent(t, users, event.ID)

	// Update leaderboard concurrently
	updateLeaderboard(t, users, event.ID)

	// Get leaderboard
	getLeaderboard(t, event.ID)
}

func createUsers(t *testing.T) []models.User {
	var wg sync.WaitGroup
	users := make([]models.User, userCount)
	userChan := make(chan models.User, userCount)

	for i := 0; i < userCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			user := createUser(t)
			userChan <- user
		}()
	}

	go func() {
		wg.Wait()
		close(userChan)
	}()

	for user := range userChan {
		users = append(users, user)
	}

	return users
}

func createUser(t *testing.T) models.User {
	username := fmt.Sprintf("testuser-%s", uuid.New().String())
	payload := map[string]string{"username": username, "country": "US"}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/user/profile", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer resp.Body.Close()

	var user models.User
	json.NewDecoder(resp.Body).Decode(&user)
	return user
}

func setUserLevels(t *testing.T, users []models.User) {
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(u models.User) {
			defer wg.Done()
			setUserLevel(t, u.ID)
		}(user)
	}
	wg.Wait()
}

func setUserLevel(t *testing.T, userID string) {
	level := rand.Intn(20) + 81 // Random level between 81 and 100
	payload := map[string]interface{}{"user_id": userID, "level": level}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/user/level", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatalf("Failed to set user level: %v", err)
	}
	defer resp.Body.Close()
}

func createEvent(t *testing.T) models.Event {
	name := fmt.Sprintf("Event-%s", time.Now().Format(time.RFC3339))
	startTime := time.Now().UTC()
	endTime := startTime.Add(30 * 24 * time.Hour)

	payload := map[string]string{
		"name":       name,
		"start_time": startTime.Format(time.RFC3339),
		"end_time":   endTime.Format(time.RFC3339),
	}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/admin/event", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}
	defer resp.Body.Close()

	var event models.Event
	json.NewDecoder(resp.Body).Decode(&event)
	return event
}

func joinUsersToEvent(t *testing.T, users []models.User, eventID string) {
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(u models.User) {
			defer wg.Done()
			joinEvent(t, eventID, u.ID)
		}(user)
	}
	wg.Wait()
}

func joinEvent(t *testing.T, eventID, userID string) {
	payload := map[string]string{"event_id": eventID, "user_id": userID}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/event/join", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatalf("Failed to join event: %v", err)
	}
	defer resp.Body.Close()
}

func updateLeaderboard(t *testing.T, users []models.User, eventID string) {
	var wg sync.WaitGroup
	for i := 0; i < updateCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			randomUser := users[rand.Intn(len(users))]
			updateLeaderboardProgress(t, eventID, randomUser.ID)
		}()
	}
	wg.Wait()
}

func updateLeaderboardProgress(t *testing.T, eventID, userID string) {
	score := rand.Intn(10000) + 1
	payload := map[string]interface{}{
		"event_id": eventID,
		"user_id":  userID,
		"score":    score,
	}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/event/leaderboard/progress", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatalf("Failed to update leaderboard: %v", err)
	}
	defer resp.Body.Close()
}

func getLeaderboard(t *testing.T, eventID string) {
	url := fmt.Sprintf("%s/event/leaderboard?event_id=%s&limit=100", baseURL, eventID)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to get leaderboard: %v", err)
	}
	defer resp.Body.Close()

	var leaderboard []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&leaderboard)

	if len(leaderboard) > 100 {
		t.Errorf("Leaderboard returned more than 100 entries: %d", len(leaderboard))
	}
}
