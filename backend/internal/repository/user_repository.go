package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/igorracki/f1/backend/internal/models"
)

type UserRepository struct {
	database *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{database: db}
}

func (userRepo *UserRepository) CreateUser(user *models.User, profile *models.Profile) error {
	log.Printf("INFO: Attempting to create user [username: %s, email: %s]", user.Username, user.Email)

	transaction, err := userRepo.database.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction for user creation: %w", err)
	}
	defer transaction.Rollback()

	_, err = transaction.Exec(
		"INSERT INTO users (id, username, email, created_at) VALUES (?, ?, ?, ?)",
		user.ID, user.Username, user.Email, user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting user %s: %w", user.ID, err)
	}

	_, err = transaction.Exec(
		"INSERT INTO profiles (user_id, display_name) VALUES (?, ?)",
		profile.UserID, profile.DisplayName,
	)
	if err != nil {
		return fmt.Errorf("inserting profile for user %s: %w", profile.UserID, err)
	}

	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("committing transaction for user %s: %w", user.ID, err)
	}

	log.Printf("INFO: Successfully created user [id: %s, username: %s]", user.ID, user.Username)
	return nil
}

func (userRepo *UserRepository) GetUserByID(id string) (*models.User, error) {
	log.Printf("INFO: Fetching user [id: %s]", id)

	user := &models.User{}
	err := userRepo.database.QueryRow(
		"SELECT id, username, email, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("INFO: User not found [id: %s]", id)
			return nil, nil
		}
		return nil, fmt.Errorf("querying user %s: %w", id, err)
	}

	log.Printf("INFO: Successfully fetched user [id: %s]", id)
	return user, nil
}
