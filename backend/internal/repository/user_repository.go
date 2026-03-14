package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/igorracki/f1/backend/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User, passwordHash string, profile *models.Profile) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, string, error)
	GetProfileByUserID(ctx context.Context, userID string) (*models.Profile, error)
}

type userRepository struct {
	database *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{database: db}
}

func (userRepo *userRepository) CreateUser(ctx context.Context, user *models.User, passwordHash string, profile *models.Profile) error {
	log.Printf("INFO: Attempting to create user [email: %s]", user.Email)

	transaction, err := userRepo.database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction for user creation: %w", err)
	}
	defer transaction.Rollback()

	_, err = transaction.ExecContext(ctx,
		"INSERT INTO users (id, email, password_hash, created_at) VALUES (?, ?, ?, ?)",
		user.ID, user.Email, passwordHash, user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting user %s: %w", user.ID, err)
	}

	_, err = transaction.ExecContext(ctx,
		"INSERT INTO profiles (user_id, display_name) VALUES (?, ?)",
		profile.UserID, profile.DisplayName,
	)
	if err != nil {
		return fmt.Errorf("inserting profile for user %s: %w", profile.UserID, err)
	}

	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("committing transaction for user %s: %w", user.ID, err)
	}

	log.Printf("INFO: Successfully created user [id: %s, email: %s]", user.ID, user.Email)
	return nil
}

func (userRepo *userRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	log.Printf("INFO: Fetching user [id: %s]", id)

	user := &models.User{}
	err := userRepo.database.QueryRowContext(ctx,
		"SELECT id, email, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Email, &user.CreatedAt)

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

func (userRepo *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, string, error) {
	log.Printf("INFO: Fetching user [email: %s]", email)

	user := &models.User{}
	var passwordHash string
	err := userRepo.database.QueryRowContext(ctx,
		"SELECT id, email, password_hash, created_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &passwordHash, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("INFO: User not found [email: %s]", email)
			return nil, "", nil
		}
		return nil, "", fmt.Errorf("querying user by email %s: %w", email, err)
	}

	log.Printf("INFO: Successfully fetched user [id: %s, email: %s]", user.ID, user.Email)
	return user, passwordHash, nil
}

func (userRepo *userRepository) GetProfileByUserID(ctx context.Context, userID string) (*models.Profile, error) {
	log.Printf("INFO: Fetching profile [user_id: %s]", userID)

	profile := &models.Profile{}
	err := userRepo.database.QueryRowContext(ctx,
		"SELECT user_id, display_name FROM profiles WHERE user_id = ?",
		userID,
	).Scan(&profile.UserID, &profile.DisplayName)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("INFO: Profile not found [user_id: %s]", userID)
			return nil, nil
		}
		return nil, fmt.Errorf("querying profile for user %s: %w", userID, err)
	}

	log.Printf("INFO: Successfully fetched profile [user_id: %s]", userID)
	return profile, nil
}
