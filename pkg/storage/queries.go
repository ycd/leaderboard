package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ycd/leaderboard/pkg/queries"
)

// GetLeaderboard returns the results from leaderboard table.
func (s *Storage) GetLeaderboard() (interface{}, error) {
	ctx := context.Background()
	cachedData, err := s.CacheGet(ctx, "leaderboard")
	if err == nil {
		return cachedData, nil
	}

	var data []interface{}

	rows, err := s.connection.Query(context.Background(), queries.GetLeaderboard)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var lr LeaderboardResult
		err := rows.Scan(&lr.Rank, &lr.Points, &lr.DisplayName, &lr.Country)
		if err != nil {
			log.Println("Got error:", err)
		}

		data = append(data, lr)
	}

	json, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if err := s.CacheSet(ctx, "leaderboard", json); err != nil {
		return nil, err
	}

	return data, nil
}

// GetLeaderboardWithCountry returns the results from leaderboard table with specific country.
func (s *Storage) GetLeaderboardWithCountry(country string) (interface{}, error) {
	rows, err := s.connection.Query(context.Background(), queries.GetLeaderboardWithCountry, country)
	if err != nil {
		return nil, err
	}

	var data []interface{}
	for rows.Next() {
		var lr LeaderboardResult
		err := rows.Scan(&lr.Rank, &lr.Points, &lr.DisplayName, &lr.Country)
		if err != nil {
			log.Println("Got error:", err)
		}

		data = append(data, lr)
	}

	return data, nil
}

// ScoreSubmit inserts a new score record to the database.
func (s *Storage) ScoreSubmit(ScoreWorth float64, UserID string, Timestamp int64) (interface{}, error) {
	_, err := s.connection.Exec(context.Background(), queries.InsertScore, UserID, ScoreWorth, Timestamp)
	if err != nil {
		return nil, err
	}

	return Score{
		ScoreWorth: ScoreWorth,
		UserID:     UserID,
		Timestamp:  Timestamp,
	}, nil
}

// UserCreate creates a new user.
func (s *Storage) UserCreate(userID, userName, country string) (interface{}, error) {
	_, err := s.connection.Exec(context.Background(), queries.NewUser, userID, userName, country)
	if err != nil {
		return nil, fmt.Errorf("user with name: %s already exists", userName)
	}

	user, err := s.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves the information about the user.
func (s *Storage) GetUser(userID string) (interface{}, error) {
	var u UserInfo

	err := s.connection.QueryRow(context.Background(), queries.GetUser, userID).Scan(&u.UserID, &u.DisplayName, &u.Points, &u.Rank)
	if err != nil {
		return nil, err
	}

	return u, nil
}
