package storage

type UserInfo struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Points      int64  `json:"points"`
	Rank        int64  `json:"rank"`
}

type LeaderboardResult struct {
	Rank        int64  `json:"rank"`
	Points      int64  `json:"points"`
	DisplayName string `json:"display_name"`
	Country     string `json:"country"`
}

type Score struct {
	ScoreWorth float64 `json:"score_worth"`
	UserID     string  `json:"user_id"`
	Timestamp  int64   `json:"timestamp"`
}
