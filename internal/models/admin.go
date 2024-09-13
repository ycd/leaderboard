package models

import "time"

type CreateEventRequest struct {
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type CreateEventResponse struct {
	Event Event `json:"event"`
}

type ListEventsResponse struct {
	Events []Event `json:"events"`
}

type DeleteEventRequest struct {
	EventID string `json:"event_id"`
}

type DeleteEventResponse struct {
	Success bool `json:"success"`
}

type GetUserDetailsRequest struct {
	UserID string `json:"user_id"`
}

type GetUserDetailsResponse struct {
	User User `json:"user"`
}

type ListUsersResponse struct {
	Users []User `json:"users"`
}

type BanUserRequest struct {
	UserID string `json:"user_id"`
}

type BanUserResponse struct {
	Success bool `json:"success"`
}

type UnbanUserRequest struct {
	UserID string `json:"user_id"`
}

type UnbanUserResponse struct {
	Success bool `json:"success"`
}

type ResetLeaderboardRequest struct {
	EventID string `json:"event_id"`
}

type ResetLeaderboardResponse struct {
	Success bool `json:"success"`
}

type Admin struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
}
