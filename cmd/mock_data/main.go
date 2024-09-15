package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ycd/leaderboard/internal/models"
)

const (
	baseURL = "http://api.leaderboard.yagizdegirmenci.com"
)

func sendRequest(method, endpoint string, payload interface{}) (*http.Response, error) {
	var req *http.Request
	var err error

	if payload != nil {
		jsonPayload, _ := json.Marshal(payload)
		req, err = http.NewRequest(method, baseURL+endpoint, bytes.NewBuffer(jsonPayload))
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, baseURL+endpoint, nil)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func createUser() models.User {
	username := fmt.Sprintf("testuser-%s", uuid.New().String())
	payload := map[string]string{
		"username": username,
		"country":  "US",
	}

	resp, err := sendRequest("POST", "/user/profile", payload)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Printf("Failed to decode user: %v", err)
		return models.User{}
	}
	return user
}

func main() {
	numUsers := 70
	eventMul := 15

	var users []models.User

	var wg sync.WaitGroup

	// Timing user creation
	startTime := time.Now()
	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			user := createUser()
			if user.ID == "" {
				log.Printf("Failed to create user")
				return
			}
			users = append(users, user)
		}()
	}
	wg.Wait()
	userCreationTime := time.Since(startTime)
	log.Printf("Created %d users in %v", len(users), userCreationTime)

	// Timing level updates
	startTime = time.Now()
	for _, user := range users {
		wg.Add(1)
		go func(userID string) {
			defer wg.Done()
			payload := map[string]interface{}{
				"user_id": userID,
				"level":   144,
			}

			resp, err := sendRequest("POST", "/user/level", payload)
			if err != nil {
				log.Printf("Failed to update level: %v", err)
				return
			}

			if resp.StatusCode != http.StatusOK {
				log.Printf("Failed to update level: %v", resp.StatusCode)
				return
			}

			defer resp.Body.Close()
		}(user.ID)
	}
	wg.Wait()
	levelUpdateTime := time.Since(startTime)
	log.Printf("Updated level for %d users in %v", len(users), levelUpdateTime)

	// Timing event creation
	startTime = time.Now()
	payload := map[string]string{
		"name":       "masters",
		"start_time": "2024-09-01T00:00:00Z",
		"end_time":   "2024-09-30T23:59:59Z",
	}

	resp, err := sendRequest("POST", "/admin/event", payload)
	if err != nil {
		log.Fatalf("Failed to create event: %v", err)
	}
	defer resp.Body.Close()

	var event models.CreateEventResponse
	if err := json.NewDecoder(resp.Body).Decode(&event); err != nil {
		log.Printf("Failed to decode event: %v", err)
	}
	eventCreationTime := time.Since(startTime)
	log.Printf("Created event: %v in %v", event, eventCreationTime)
	log.Printf("Event ID: %s", event.Event.ID)

	// Timing joining events
	startTime = time.Now()
	for _, user := range users {
		if user.ID == "" {
			continue
		}
		wg.Add(1)
		go func(userID string) {
			defer wg.Done()
			payload := map[string]interface{}{
				"user_id":  userID,
				"event_id": event.Event.ID,
			}

			resp, err := sendRequest("POST", "/event/join", payload)
			if err != nil {
				return
			}
			if resp.StatusCode != http.StatusOK {
				return
			}

			defer resp.Body.Close()
		}(user.ID)

		wg.Wait()
	}
	joinEventTime := time.Since(startTime)
	log.Printf("Joined event for %d users in %v", len(users), joinEventTime)

	// Timing leaderboard updates
	startTime = time.Now()
	for i := 0; i < eventMul; i++ {
		for _, user := range users {
			wg.Add(1)
			go func(userID string) {
				defer wg.Done()
				payload := map[string]interface{}{
					"user_id":  userID,
					"event_id": event.Event.ID,
					"score":    rand.Intn(100),
				}

				resp, err := sendRequest("POST", "/event/leaderboard/progress", payload)
				if err != nil {
					log.Printf("Failed to update leaderboard: %v", err)
					return
				}

				if resp.StatusCode != http.StatusOK {
					log.Printf("Failed to update leaderboard: %v", resp.StatusCode)
					return
				}

				defer resp.Body.Close()
			}(user.ID)
		}
	}
	wg.Wait()
	leaderboardUpdateTime := time.Since(startTime)
	log.Printf("Updated leaderboard with %d progress events in %v", len(users)*10, leaderboardUpdateTime)

	// Timing leaderboard retrieval
	startTime = time.Now()
	ldb, err := http.Get(baseURL + "/event/leaderboard?event_id=" + event.Event.ID + "&limit=100")
	if err != nil {
		log.Printf("Failed to get leaderboard: %v", err)
	}
	defer ldb.Body.Close()

	var leaderboard models.LeaderboardResponse
	if err := json.NewDecoder(ldb.Body).Decode(&leaderboard); err != nil {
		log.Printf("Failed to decode leaderboard: %v", err)
	}
	leaderboardRetrievalTime := time.Since(startTime)
	log.Printf("Retrieved leaderboard in %v", leaderboardRetrievalTime)

	log.Printf("Load test completed")
	log.Printf("Total users created: %d", len(users))
	log.Printf("Total events created: 1")
	log.Printf("Total leaderboard updates: %d", len(leaderboard.Entries)*eventMul)
}
