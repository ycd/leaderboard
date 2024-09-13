package repository

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/ycd/leaderboard/internal/errors"
	"github.com/ycd/leaderboard/internal/models"
)

type UserRepository struct {
	db *sql.DB // This should be the SQLite connection
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) (string, error) {
	newUUID := uuid.New().String()
	query := `INSERT INTO users (id, username, country, level, coin) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, newUUID, user.Username, user.Country, user.Level, user.Coin)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return "", err
	}
	log.Printf("User created with ID: %s", newUUID)
	return newUUID, nil // Correctly return the UUID
}

func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, country, level, coin, is_banned FROM users WHERE id = ?"
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Country, &user.Level, &user.Coin, &user.IsBanned)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User not found with ID: %s", id)
			return nil, errors.NewError(errors.ErrUserNotFound, "user not found")
		}
		log.Printf("Error retrieving user by ID: %v", err)
		return nil, err
	}
	log.Printf("User retrieved with ID: %s", id)
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, country, level, coin FROM users WHERE username = ?`
	var user models.User
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Country, &user.Level, &user.Coin)
	if err == sql.ErrNoRows {
		log.Printf("User not found with username: %s", username)
		return nil, nil
	}
	if err != nil {
		log.Printf("Error retrieving user by username: %v", err)
		return nil, err
	}
	log.Printf("User retrieved with username: %s", username)
	return &user, nil
}

func (r *UserRepository) SetUserLevel(userID string, level int) error {
	query := `UPDATE users SET level = ? WHERE id = ?`
	_, err := r.db.Exec(query, level, userID)
	if err != nil {
		log.Printf("Error setting user level for userID %s: %v", userID, err)
	}
	log.Printf("User level set to %d for userID %s", level, userID)
	return err
}

func (r *UserRepository) SetUserCoin(userID string, coin int) error {
	query := `UPDATE users SET coin = ? WHERE id = ?`
	_, err := r.db.Exec(query, coin, userID)
	if err != nil {
		log.Printf("Error setting user coin for userID %s: %v", userID, err)
	}
	log.Printf("User coin set to %d for userID %s", coin, userID)
	return err
}

func (r *UserRepository) AddUserCoin(userID string, coin int) error {
	query := `UPDATE users SET coin = coin + ? WHERE id = ?`
	_, err := r.db.Exec(query, coin, userID)
	if err != nil {
		log.Printf("Error adding user coin for userID %s: %v", userID, err)
	}
	log.Printf("Added %d coins to userID %s", coin, userID)
	return err
}
