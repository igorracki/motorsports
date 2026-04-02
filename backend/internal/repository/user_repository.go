package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/igorracki/motorsports/backend/internal/database"
	"github.com/igorracki/motorsports/backend/internal/models"
)

var (
	ErrDuplicateEmail = errors.New("email already registered")
)

const (
	insertUserSQL         = "INSERT INTO users (id, email, password_hash, created_at) VALUES (?, ?, ?, ?)"
	insertProfileSQL      = "INSERT INTO profiles (user_id, display_name) VALUES (?, ?)"
	getUserByIDSQL        = "SELECT id, email, created_at FROM users WHERE id = ?"
	getUserByEmailSQL     = "SELECT id, email, password_hash, created_at FROM users WHERE email = ?"
	getProfileByUserIDSQL = "SELECT user_id, display_name FROM profiles WHERE user_id = ?"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User, passwordHash string, profile *models.Profile) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, string, error)
	GetProfileByUserID(ctx context.Context, userID string) (*models.Profile, error)
	GetProfilesByUserIDs(ctx context.Context, userIDs []string) ([]models.Profile, error)
}

type userRepository struct {
	manager *database.Manager
}

func NewUserRepository(manager *database.Manager) UserRepository {
	return &userRepository{manager: manager}
}

func (repo *userRepository) CreateUser(ctx context.Context, user *models.User, passwordHash string, profile *models.Profile) error {
	return repo.manager.Transaction(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, insertUserSQL, user.ID, user.Email, passwordHash, user.CreatedAt); err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
				return ErrDuplicateEmail
			}
			return fmt.Errorf("inserting user %s: %w", user.ID, err)
		}

		if _, err := tx.ExecContext(ctx, insertProfileSQL, profile.UserID, profile.DisplayName); err != nil {
			return fmt.Errorf("inserting profile for user %s: %w", profile.UserID, err)
		}
		return nil
	})
}

func (repo *userRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	err := repo.manager.DB().QueryRowContext(ctx, getUserByIDSQL, id).Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("querying user %s: %w", id, err)
	}

	return user, nil
}

func (repo *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, string, error) {
	user := &models.User{}
	var passwordHash string
	err := repo.manager.DB().QueryRowContext(ctx, getUserByEmailSQL, email).Scan(&user.ID, &user.Email, &passwordHash, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", nil
		}
		return nil, "", fmt.Errorf("querying user by email: %w", err)
	}

	return user, passwordHash, nil
}

func (repo *userRepository) GetProfileByUserID(ctx context.Context, userID string) (*models.Profile, error) {
	profile := &models.Profile{}
	err := repo.manager.DB().QueryRowContext(ctx, getProfileByUserIDSQL, userID).Scan(&profile.UserID, &profile.DisplayName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("querying profile for user %s: %w", userID, err)
	}

	return profile, nil
}

func (repo *userRepository) GetProfilesByUserIDs(ctx context.Context, userIDs []string) ([]models.Profile, error) {
	if len(userIDs) == 0 {
		return []models.Profile{}, nil
	}

	placeholders := database.GeneratePlaceholders(len(userIDs))
	query := fmt.Sprintf("SELECT user_id, display_name FROM profiles WHERE user_id IN (%s)", placeholders)

	rows, err := repo.manager.DB().QueryContext(ctx, query, database.ToAnySlice(userIDs)...)
	if err != nil {
		return nil, fmt.Errorf("querying profiles: %w", err)
	}
	defer rows.Close()

	profiles := make([]models.Profile, 0)
	for rows.Next() {
		var profile models.Profile
		if err := rows.Scan(&profile.UserID, &profile.DisplayName); err != nil {
			return nil, fmt.Errorf("scanning profile: %w", err)
		}
		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating profiles: %w", err)
	}

	return profiles, nil
}
