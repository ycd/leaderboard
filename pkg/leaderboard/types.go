package leaderboard

type ScoreSubmit struct {
	ScoreWorth float32 `json:"score_worth"`
	UserID     string  `json:"user_id"`
	Timestamp  int     `json:"timestamp"`
}
