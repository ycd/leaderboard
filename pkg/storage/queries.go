package storage

import (
	"context"
	"fmt"
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

		data = append(data, v)
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
func (s *Storage) ScoreSubmit(ScoreWorth float32, UserID string, Timestamp int64) (interface{}, error) {
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

// UserCreate creates a new user.
func (s *Storage) UserCreate(userID, userName, country string) (interface{}, error) {
	_, err := s.connection.Query(context.Background(), queries.NewUser, userID, userName, country)
	if err != nil {
		return nil, fmt.Errorf("user with name: %s already exists", userName)
	}

	user, err := s.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) GetUser(userID string) (interface{}, error) {
	var u UserInfo

	row := s.connection.QueryRow(context.Background(), queries.GetUser, userID)
	err := row.Scan(&u.UserID, &u.DisplayName, &u.Country, &u.Points, &u.Rank)
	// if err != nil {
	// 	if err, ok := err.(*pgx.Errro)
	// }

	log.Println(err.Error())

	return u, nil
}
