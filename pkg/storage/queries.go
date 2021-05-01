package storage

import (
	"context"
	"log"

	"github.com/ycd/leaderboard/pkg/queries"
)

// GetLeaderboard returns the results from leaderboard table.
func (s *Storage) GetLeaderboard() (interface{}, error) {
	rows, err := s.connection.Query(context.Background(), queries.GetLeaderboard)
	if err != nil {
		return nil, err
	}

	var data []interface{}
	for rows.Next() {
		v, err := rows.Values()
		if err != nil {
			log.Println("Got error:", err)
		}

		log.Println("Got value:", v)
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
		v, err := rows.Values()
		if err != nil {
			log.Println("Got error:", err)
		}

		data = append(data, v)
	}

	return data, nil
}

// ScoreSubmit inserts a new score record to the database.
func (s *Storage) ScoreSubmit(ScoreWorth float32, UserID string, Timestamp int) (interface{}, error) {
	rows, err := s.connection.Query(context.Background(), queries.InsertScore, UserID, ScoreWorth, Timestamp)
	if err != nil {
		return nil, err
	}

	var data []interface{}
	for rows.Next() {
		v, err := rows.Values()
		if err != nil {
			log.Println("Got error:", err)
		}

		data = append(data, v)
	}

	return data, nil
}
