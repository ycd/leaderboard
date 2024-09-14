package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/ycd/leaderboard/internal/errors"
	"github.com/ycd/leaderboard/internal/models"
)

type EventRepository struct {
	chDB     *sql.DB
	sqliteDB *sql.DB
	cacheDB  *redis.Client
}

func NewEventRepository(chDB *sql.DB, sqliteDB *sql.DB, cacheDB *redis.Client) *EventRepository {
	return &EventRepository{
		chDB:     chDB,
		sqliteDB: sqliteDB,
		cacheDB:  cacheDB,
	}
}

func (r *EventRepository) checkEventExists(eventID string) error {
	query := `SELECT id, name FROM events WHERE id = ?`
	var event models.Event
	err := r.chDB.QueryRow(query, eventID).Scan(&event.ID, &event.Name)
	if err == sql.ErrNoRows {
		return errors.NewError(errors.ErrEventNotFound, "Event not found")
	}

	if err != nil {
		log.Printf("Error checking if event %s exists: %v", eventID, err)
		return errors.NewError(errors.ErrEventNotFound, "Event not found")
	}
	log.Printf("Event %s exists.", eventID)
	return nil
}

func (r *EventRepository) ClaimReward(eventID, userID string, rewardAmount int) error {
	query := `
    INSERT OR IGNORE INTO rewards (event_id, user_id, reward_amount)
    VALUES (?, ?, ?)
    `
	result, err := r.sqliteDB.Exec(query, eventID, userID, rewardAmount)
	if err != nil {
		return fmt.Errorf("error claiming reward: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewError(errors.ErrRewardAlreadyClaimed, "Reward already claimed")
	}

	return nil
}

func (r *EventRepository) getParticipantCountry(eventID, userID string) (string, error) {
	cacheKey := fmt.Sprintf("event:%s:user:%s:country", eventID, userID)
	log.Printf("GetParticipantCountry cacheKey: %s", cacheKey)

	country, err := r.cacheDB.Get(context.Background(), cacheKey).Result()
	if err == redis.Nil {
		var countryResult string
		query := `SELECT country FROM users WHERE id = ?`
		err := r.sqliteDB.QueryRow(query, userID).Scan(&countryResult)
		if err != nil {
			log.Printf("Error retrieving country for user %s in event %s: %v", userID, eventID, err)
			return "", errors.NewError(errors.ErrDatabaseError, "Error retrieving country")
		}
		err = r.cacheDB.Set(context.Background(), cacheKey, countryResult, 0).Err()
		if err != nil {
			log.Printf("Error setting cache for user %s in event %s: %v", userID, eventID, err)
		}
		log.Printf("Participant country: %s", countryResult)
		return countryResult, nil
	} else if err != nil {
		log.Printf("Error retrieving country from cache for user %s in event %s: %v", userID, eventID, err)
		return "", errors.NewError(errors.ErrDatabaseError, "Error retrieving country from cache")
	}

	return country, nil
}

func (r *EventRepository) AddParticipant(eventID, userID string) error {
	query := `
    INSERT INTO event_participants (event_id, user_id, score, country)
    VALUES (?, ?, ?, ?)
    `

	country, err := r.getParticipantCountry(eventID, userID)
	if err != nil {
		return fmt.Errorf("error getting participant country: %w", err)
	}

	log.Printf("Executing query: %s with params: event_id=%s, user_id=%s, score=0, country=%s", query, eventID, userID, country)

	_, err = r.chDB.Exec(query, eventID, userID, 0, country)
	if err != nil {
		return fmt.Errorf("error inserting participant: %w", err)
	}

	log.Printf("Successfully added participant %s to event %s", userID, eventID)
	return nil
}

func (r *EventRepository) AddParticipantScore(eventID, userID string, score int) error {
	query := `
    INSERT INTO event_participants (event_id, user_id, score, country)
    VALUES (?, ?, ?, ?)
    `

	country, err := r.getParticipantCountry(eventID, userID)
	if err != nil {
		return fmt.Errorf("error getting participant country: %w", err)
	}

	log.Printf("Executing query: %s with params: event_id=%s, user_id=%s, score=%d, country=%s", query, eventID, userID, score, country)

	_, err = r.chDB.Exec(query, eventID, userID, score, country)
	if err != nil {
		return fmt.Errorf("error updating participant score: %w", err)
	}

	return nil
}

func (r *EventRepository) isUserJoined(eventID, userID string) (bool, error) {
	query := `SELECT COUNT(*) FROM event_participants WHERE event_id = ? AND user_id = ?`
	var count int
	err := r.chDB.QueryRow(query, eventID, userID).Scan(&count)
	if err != nil {
		log.Printf("Error checking if user %s has joined event %s: %v", userID, eventID, err)
		return false, errors.NewError(errors.ErrDatabaseError, "Error checking user participation")
	}
	log.Printf("User %s has joined event %s: %t", userID, eventID, count > 0)
	return count > 0, nil
}

func (r *EventRepository) GetLeaderboard(eventID string, limit int) ([]models.LeaderboardEntryWithRank, error) {
	query := `
SELECT
    user_id,
    SUM(score) AS total_score,
    RANK() OVER (ORDER BY SUM(score) DESC) AS rank
FROM event_participants
WHERE event_id = ?
GROUP BY user_id
ORDER BY total_score DESC
LIMIT ?
    `

	rows, err := r.chDB.Query(query, eventID, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying leaderboard: %w", err)
	}
	defer rows.Close()

	var leaderboard []models.LeaderboardEntryWithRank
	for rows.Next() {
		var entry models.LeaderboardEntryWithRank
		if err := rows.Scan(&entry.UserID, &entry.TotalScore, &entry.Rank); err != nil {
			return nil, fmt.Errorf("error scanning leaderboard entry: %w", err)
		}
		leaderboard = append(leaderboard, entry)
	}

	return leaderboard, nil
}

func (r *EventRepository) GetCountryLeaderboard(eventID string, country string, limit int) ([]models.CountryLeaderboardEntryWithRank, error) {
	query := `
SELECT
    user_id,
	country,
    SUM(score) AS total_score,
    RANK() OVER (ORDER BY SUM(score) DESC) AS rank
FROM event_participants
WHERE event_id = ?
AND country = ?
GROUP BY user_id, country
ORDER BY total_score DESC
LIMIT ?
    `

	rows, err := r.chDB.Query(query, eventID, country, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying country leaderboard: %w", err)
	}
	defer rows.Close()

	var leaderboard []models.CountryLeaderboardEntryWithRank
	for rows.Next() {
		var entry models.CountryLeaderboardEntryWithRank
		if err := rows.Scan(&entry.UserID, &entry.Country, &entry.TotalScore, &entry.Rank); err != nil {
			return nil, fmt.Errorf("error scanning country leaderboard entry: %w", err)
		}
		leaderboard = append(leaderboard, entry)
	}

	return leaderboard, nil
}

func (r *EventRepository) GetParticipantRank(eventID, userID string) (int, error) {
	query := `
        SELECT rank FROM (
            SELECT
                user_id,
                RANK() OVER (ORDER BY total_score DESC) as rank
            FROM event_leaderboard_mv
            WHERE event_id = ?
        ) ranked
        WHERE user_id = ?
    `
	var rank int
	err := r.chDB.QueryRow(query, eventID, userID).Scan(&rank)
	if err != nil {
		log.Printf("Error retrieving rank for user %s in event %s: %v", userID, eventID, err)
		return 0, err
	}
	log.Printf("Retrieved rank %d for user %s in event %s.", rank, userID, eventID)
	return rank, nil
}

func (r *EventRepository) CreateEvent(ctx context.Context, event *models.Event) error {
	query := `
		INSERT INTO events (id, name, start_time, end_time)
		VALUES (?, ?, ?, ?)
	`
	_, err := r.sqliteDB.ExecContext(ctx, query, event.ID, event.Name, event.StartTime, event.EndTime)
	return err
}

func (r *EventRepository) GetEventsByUserID(ctx context.Context, userID string, limit int) ([]*models.Event, error) {
	query := `
		SELECT id, name, start_time, end_time
		FROM events
		WHERE user_id = ?

		LIMIT ?
	`
	rows, err := r.chDB.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Name, &event.StartTime, &event.EndTime)
		if err != nil {
			return nil, errors.NewError(errors.ErrEventNotFound, "Event not found")
		}
		events = append(events, &event)
	}
	return events, nil
}

func (r *EventRepository) GetEventByID(ctx context.Context, eventID string) (models.Event, error) {
	query := `
		SELECT id, name, start_time, end_time
		FROM events
		WHERE id = ?
	`

	var event models.Event
	err := r.chDB.QueryRowContext(ctx, query, eventID).Scan(&event.ID, &event.Name, &event.StartTime, &event.EndTime)
	if err == sql.ErrNoRows {
		return event, errors.NewError(errors.ErrEventNotFound, "Event not found")
	}
	if err != nil {
		return event, err
	}

	return event, nil
}
