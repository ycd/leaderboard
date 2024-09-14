package models

// SetProfileRequest represents the request body for setting a user's profile
type SetProfileRequest struct {
	Username string `json:"username"`
	Country  string `json:"country"`
	Level    int    `json:"level"`
	Coin     int    `json:"coin"`
}

// SetLevelRequest represents the request body for setting a user's level
type SetLevelRequest struct {
	UserID string `json:"user_id"`
	Level  int    `json:"level"`
}

// SetCoinRequest represents the request body for setting a user's coin amount
type SetCoinRequest struct {
	UserID string `json:"user_id"`
	Coin   int    `json:"coin"`
}

// JoinEventRequest represents the request body for joining an event
type JoinEventRequest struct {
	UserID  string `json:"user_id"`
	EventID string `json:"event_id"`
}

// SetLeaderboardProgressRequest represents the request body for updating a user's leaderboard progress
type SetLeaderboardProgressRequest struct {
	UserID  string `json:"user_id"`
	EventID string `json:"event_id"`
	Score   int    `json:"score"`
}

// GetLeaderboardRequest represents the query parameters for getting a leaderboard
type GetLeaderboardRequest struct {
	EventID string `json:"event_id"`
	Limit   int    `json:"limit"`
}

// GetCountryLeaderboardRequest represents the query parameters for getting a country leaderboard
type GetCountryLeaderboardRequest struct {
	EventID string `json:"event_id"`
	Country string `json:"country"`
	Limit   int    `json:"limit"`
}

// ClaimLeaderboardRequest represents the request body for claiming a leaderboard reward
type ClaimLeaderboardRequest struct {
	UserID  string `json:"user_id"`
	EventID string `json:"event_id"`
}
