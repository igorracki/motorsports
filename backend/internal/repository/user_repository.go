package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/igorracki/f1/backend/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User, passwordHash string, profile *models.Profile) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, string, error)
	GetProfileByUserID(ctx context.Context, userID string) (*models.Profile, error)
	GetProfilesByUserIDs(ctx context.Context, userIDs []string) ([]models.Profile, error)
}

type userRepository struct {
	database *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{database: db}
}

func (userRepo *userRepository) CreateUser(ctx context.Context, user *models.User, passwordHash string, profile *models.Profile) error {
	slog.InfoContext(ctx, "Entry: CreateUser", "email", user.Email)

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

	slog.InfoContext(ctx, "Exit: CreateUser", "user_id", user.ID, "email", user.Email)
	return nil
}

func (userRepo *userRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	slog.InfoContext(ctx, "Entry: GetUserByID", "user_id", id)

	user := &models.User{}
	err := userRepo.database.QueryRowContext(ctx,
		"SELECT id, email, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.InfoContext(ctx, "User not found", "user_id", id)
			return nil, nil
		}
		return nil, fmt.Errorf("querying user %s: %w", id, err)
	}

	slog.InfoContext(ctx, "Exit: GetUserByID", "user_id", id)
	return user, nil
}

func (userRepo *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, string, error) {
	slog.InfoContext(ctx, "Entry: GetUserByEmail", "email", email)

	user := &models.User{}
	var passwordHash string
	err := userRepo.database.QueryRowContext(ctx,
		"SELECT id, email, password_hash, created_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &passwordHash, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.InfoContext(ctx, "User not found", "email", email)
			return nil, "", nil
		}
		return nil, "", fmt.Errorf("querying user by email %s: %w", email, err)
	}

	slog.InfoContext(ctx, "Exit: GetUserByEmail", "user_id", user.ID, "email", user.Email)
	return user, passwordHash, nil
}

func (userRepo *userRepository) GetProfileByUserID(ctx context.Context, userID string) (*models.Profile, error) {
	slog.InfoContext(ctx, "Entry: GetProfileByUserID", "user_id", userID)

	profile := &models.Profile{}
	err := userRepo.database.QueryRowContext(ctx,
		"SELECT user_id, display_name FROM profiles WHERE user_id = ?",
		userID,
	).Scan(&profile.UserID, &profile.DisplayName)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.InfoContext(ctx, "Profile not found", "user_id", userID)
			return nil, nil
		}
		return nil, fmt.Errorf("querying profile for user %s: %w", userID, err)
	}

	slog.InfoContext(ctx, "Exit: GetProfileByUserID", "user_id", userID)
	return profile, nil
}

func (userRepo *userRepository) GetProfilesByUserIDs(ctx context.Context, userIDs []string) ([]models.Profile, error) {
	slog.InfoContext(ctx, "Entry: GetProfilesByUserIDs", "count", len(userIDs))

	if len(userIDs) == 0 {
		return []models.Profile{}, nil
	}

	query := "SELECT user_id, display_name FROM profiles WHERE user_id IN ("
	args := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		if i > 0 {
			query += ","
		}
		query += "?"
		args[i] = id
	}
	query += ")"

	rows, err := userRepo.database.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying profiles: %w", err)
	}
	defer rows.Close()

	profiles := []models.Profile{}
	for rows.Next() {
		var p models.Profile
		if err := rows.Scan(&p.UserID, &p.DisplayName); err != nil {
			return nil, fmt.Errorf("scanning profile: %w", err)
		}
		profiles = append(profiles, p)
	}

	slog.InfoContext(ctx, "Exit: GetProfilesByUserIDs", "found", len(profiles))
	return profiles, nil
}
