package storage

type UserInfo struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Points      int64  `json:"points,omitempty"`
	Country     string `json:"country"`
	Rank        int64  `json:"rank"`
}
