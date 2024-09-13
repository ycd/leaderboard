package errors

import (
	"fmt"
)

type ErrorType string

const (
	ErrUserNotFound               ErrorType = "user not found"
	ErrUserAlreadyExists          ErrorType = "user already exists"
	ErrEventNotFound              ErrorType = "event not found"
	ErrInvalidUsername            ErrorType = "invalid username"
	ErrInvalidCountry             ErrorType = "invalid country"
	ErrInvalidLevel               ErrorType = "invalid level"
	ErrInvalidCoin                ErrorType = "invalid coin amount"
	ErrInsufficientLevel          ErrorType = "insufficient level to join event"
	ErrUserAlreadyJoined          ErrorType = "user has already joined this event"
	ErrEventEnded                 ErrorType = "event has already ended"
	ErrEventNotStarted            ErrorType = "event has not started yet"
	ErrInvalidScore               ErrorType = "invalid score"
	ErrNotEligibleForReward       ErrorType = "user is not eligible for reward"
	ErrRewardAlreadyClaimed       ErrorType = "reward has already been claimed"
	ErrInvalidEventDates          ErrorType = "invalid event start or end dates"
	ErrUnauthorized               ErrorType = "unauthorized access"
	ErrInternalServerError        ErrorType = "internal server error"
	ErrDatabaseError              ErrorType = "database error"
	ErrInvalidRequest             ErrorType = "invalid request"
	ErrUserBanned                 ErrorType = "user is banned"
	ErrLeaderboardNotFound        ErrorType = "leaderboard not found"
	ErrCountryLeaderboardNotFound ErrorType = "country leaderboard not found"
)

type Error struct {
	Type    ErrorType
	Message string
}

func (e Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %s", e.Type, e.Message)
	}
	return string(e.Type)
}

func (e Error) Is(target error) bool {
	t, ok := target.(Error)
	if !ok {
		return false
	}
	return e.Type == t.Type
}

func NewError(errType ErrorType, message string) Error {
	return Error{Type: errType, Message: message}
}
