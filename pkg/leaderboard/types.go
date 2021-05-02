package leaderboard

type ScoreSubmit struct {
	ScoreWorth float32 `json:"score_worth"`
	UserID     string  `json:"user_id"`
}

type UserCreate struct {
	UserName string `json:"display_name"`
	Country  string `json:"country"`
}
