package models

import "time"

type Event struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type LeaderboardEntry struct {
	EventID  string `json:"event_id"`
	UserID   string `json:"user_id"`
	Score    int32  `json:"score"`
	Country  string `json:"country"`
	Username string `json:"username"`
	Rank     int    `json:"rank"`
}

type EventParticipant struct {
	ID      string `json:"id"`
	EventID string `json:"event_id"`
	UserID  string `json:"user_id"`
	Score   int32  `json:"score"`
}

type LeaderboardEntryWithRank struct {
	UserID     string `json:"user_id"`
	TotalScore int32  `json:"total_score"`
	Rank       int    `json:"rank"`
}

type CountryLeaderboardEntryWithRank struct {
	UserID     string `json:"user_id"`
	Country    string `json:"country"`
	TotalScore int32  `json:"total_score"`
	Rank       int    `json:"rank"`
}
