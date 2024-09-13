package repository

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/ycd/leaderboard/internal/models"
)

type AdminRepository struct {
	chDB     *sql.DB
	sqliteDB *sql.DB
}

func NewAdminRepository(chDB *sql.DB, sqliteDB *sql.DB) *AdminRepository {
	return &AdminRepository{
		chDB:     chDB,
		sqliteDB: sqliteDB,
	}
}

func (r *AdminRepository) CreateEvent(event *models.Event) error {
	event.ID = uuid.New().String()
	query := `
        INSERT INTO events (id, name, start_time, end_time)
        VALUES (?, ?, ?, ?)
    `
	_, err := r.chDB.Exec(query, event.ID, event.Name, event.StartTime, event.EndTime)
	if err != nil {
		log.Printf("Error creating event: %v", err)
		return err
	}
	log.Printf("Event created with ID: %s", event.ID)
	return nil
}

func (r *AdminRepository) ListEvents() ([]models.Event, error) {
	query := `
		SELECT id, name, start_time, end_time
		FROM events
		ORDER BY start_time DESC
	`
	rows, err := r.chDB.Query(query)
	if err != nil {
		log.Printf("Error listing events: %v", err)
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Name, &event.StartTime, &event.EndTime)
		if err != nil {
			log.Printf("Error scanning event: %v", err)
			return nil, err
		}
		events = append(events, event)
	}

	log.Printf("Listed %d events", len(events))
	return events, nil
}

func (r *AdminRepository) DeleteEvent(eventID string) error {
	query := `
		ALTER TABLE events
		DELETE WHERE id = ?
	`
	_, err := r.chDB.Exec(query, eventID)
	if err != nil {
		log.Printf("Error deleting event with ID %s: %v", eventID, err)
		return err
	}
	log.Printf("Event with ID %s deleted", eventID)
	return nil
}

func (r *AdminRepository) GetUserDetails(userID string) (*models.User, error) {
	query := `
		SELECT id, username, country, level, coin, is_banned
		FROM users
		WHERE id = ?
	`
	var user models.User
	err := r.sqliteDB.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.Country, &user.Level, &user.Coin, &user.IsBanned,
	)
	if err != nil {
		log.Printf("Error retrieving user details for userID %s: %v", userID, err)
		return nil, err
	}
	log.Printf("Retrieved user details for userID %s", userID)
	return &user, nil
}

func (r *AdminRepository) ListUsers() ([]models.User, error) {
	query := `
		SELECT id, username, country, level, coin, is_banned
		FROM users
		ORDER BY id
	`
	rows, err := r.sqliteDB.Query(query)
	if err != nil {
		log.Printf("Error listing users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Country, &user.Level, &user.Coin, &user.IsBanned)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	log.Printf("Listed %d users", len(users))
	return users, nil
}

func (r *AdminRepository) BanUser(userID string) error {
	query := `
		UPDATE users
		SET is_banned = 1
		WHERE id = ?
	`
	_, err := r.sqliteDB.Exec(query, userID)
	if err != nil {
		log.Printf("Error banning user with ID %s: %v", userID, err)
		return err
	}
	log.Printf("User with ID %s banned", userID)
	return nil
}

func (r *AdminRepository) UnbanUser(userID string) error {
	query := `
		UPDATE users
		SET is_banned = 0
		WHERE id = ?
	`
	_, err := r.sqliteDB.Exec(query, userID)
	if err != nil {
		log.Printf("Error unbanning user with ID %s: %v", userID, err)
		return err
	}
	log.Printf("User with ID %s unbanned", userID)
	return nil
}

func (r *AdminRepository) ResetLeaderboard(eventID string) error {
	query := `
		ALTER TABLE event_participants
		DELETE WHERE event_id = ?
	`
	_, err := r.chDB.Exec(query, eventID)
	if err != nil {
		log.Printf("Error resetting leaderboard for eventID %s: %v", eventID, err)
		return err
	}
	log.Printf("Leaderboard reset for eventID %s", eventID)
	return nil
}
