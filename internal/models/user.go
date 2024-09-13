package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Country  string `json:"country"`
	Level    int32  `json:"level"`
	Coin     int32  `json:"coin"`
	IsBanned bool   `json:"is_banned"`
}
