package models

// UserProfileResponse represents the response for user profile operations
type UserProfileResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Country  string `json:"country"`
	Level    int    `json:"level"`
	Coin     int    `json:"coin"`
}

// SetLevelResponse represents the response for setting a user's level
type SetLevelResponse struct {
	UserID string `json:"user_id"`
	Level  int    `json:"level"`
}

// SetCoinResponse represents the response for setting a user's coin amount
type SetCoinResponse struct {
	UserID string `json:"user_id"`
	Coin   int    `json:"coin"`
}

// JoinEventResponse represents the response for joining an event
type JoinEventResponse struct {
	UserID  string `json:"user_id"`
	EventID string `json:"event_id"`
	Success bool   `json:"success"`
}

// SetLeaderboardProgressResponse represents the response for updating a user's leaderboard progress
type SetLeaderboardProgressResponse struct {
	UserID  string `json:"user_id"`
	EventID string `json:"event_id"`
	Score   int    `json:"score"`
	Success bool   `json:"success"`
}

// LeaderboardResponse represents the response for getting a leaderboard
type LeaderboardResponse struct {
	EventID string                     `json:"event_id"`
	Entries []LeaderboardEntryWithRank `json:"entries"`
}

// CountryLeaderboardResponse represents the response for getting a country leaderboard
type CountryLeaderboardResponse struct {
	EventID string                            `json:"event_id"`
	Country string                            `json:"country"`
	Entries []CountryLeaderboardEntryWithRank `json:"entries"`
}

// ClaimRewardResponse represents the response for claiming a leaderboard reward
type ClaimRewardResponse struct {
	UserID     string `json:"user_id"`
	EventID    string `json:"event_id"`
	Rank       int    `json:"rank"`
	RewardCoin int    `json:"reward_coin"`
	Success    bool   `json:"success"`
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Error string `json:"error"`
}
