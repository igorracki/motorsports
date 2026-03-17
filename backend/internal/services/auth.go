package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/igorracki/f1/backend/internal/auth"
	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/repository"
	"github.com/microcosm-cc/bluemonday"
)

var sanitizer = bluemonday.StrictPolicy()

type AuthService interface {
	Register(ctx context.Context, request models.RegisterUserRequest) (*models.User, *models.Profile, string, time.Time, error)
	Login(ctx context.Context, request models.LoginRequest) (*models.User, *models.Profile, string, time.Time, error)
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepo,
	}
}

func (service *authService) Register(ctx context.Context, request models.RegisterUserRequest) (*models.User, *models.Profile, string, time.Time, error) {
	slog.InfoContext(ctx, "Registering user", "email", request.Email)

	if err := service.validateRegistrationRequest(request); err != nil {
		return nil, nil, "", time.Time{}, err
	}

	passwordHash, err := auth.HashPassword(request.Password)
	if err != nil {
		return nil, nil, "", time.Time{}, fmt.Errorf("hashing password: %w", err)
	}

	userID := uuid.New().String()
	user := &models.User{
		ID:        userID,
		Email:     request.Email,
		CreatedAt: time.Now().UTC(),
	}

	displayName := sanitizer.Sanitize(request.DisplayName)
	if strings.TrimSpace(displayName) == "" {
		displayName = strings.Split(request.Email, "@")[0]
	}

	profile := &models.Profile{
		UserID:      userID,
		DisplayName: displayName,
	}

	if err := service.userRepository.CreateUser(ctx, user, passwordHash, profile); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			if strings.Contains(err.Error(), "users.email") {
				return nil, nil, "", time.Time{}, fmt.Errorf("email '%s' is already registered", request.Email)
			}
		}
		return nil, nil, "", time.Time{}, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate token for immediate authentication after registration
	expiresAt := time.Now().Add(24 * time.Hour)
	token, err := auth.GenerateToken(userID, expiresAt)
	if err != nil {
		return user, profile, "", time.Time{}, fmt.Errorf("generating token: %w", err)
	}

	slog.InfoContext(ctx, "User registered successfully", "user_id", userID)
	return user, profile, token, expiresAt, nil
}

func (service *authService) Login(ctx context.Context, request models.LoginRequest) (*models.User, *models.Profile, string, time.Time, error) {
	slog.InfoContext(ctx, "Logging in user", "email", request.Email)

	user, passwordHash, err := service.userRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, nil, "", time.Time{}, fmt.Errorf("fetching user: %w", err)
	}

	if user == nil {
		return nil, nil, "", time.Time{}, fmt.Errorf("invalid email or password")
	}

	match, err := auth.VerifyPassword(request.Password, passwordHash)
	if err != nil {
		return nil, nil, "", time.Time{}, fmt.Errorf("verifying password: %w", err)
	}

	if !match {
		return nil, nil, "", time.Time{}, fmt.Errorf("invalid email or password")
	}

	profile, err := service.userRepository.GetProfileByUserID(ctx, user.ID)
	if err != nil {
		return nil, nil, "", time.Time{}, fmt.Errorf("fetching profile: %w", err)
	}

	duration := 24 * time.Hour
	if request.RememberMe {
		duration = 30 * 24 * time.Hour
	}
	expiresAt := time.Now().Add(duration)

	token, err := auth.GenerateToken(user.ID, expiresAt)
	if err != nil {
		return nil, nil, "", time.Time{}, fmt.Errorf("generating token: %w", err)
	}

	slog.InfoContext(ctx, "User logged in successfully", "user_id", user.ID)
	return user, profile, token, expiresAt, nil
}

func (service *authService) validateRegistrationRequest(request models.RegisterUserRequest) error {
	if _, err := mail.ParseAddress(request.Email); err != nil {
		return fmt.Errorf("invalid email address format")
	}

	if len(request.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if strings.TrimSpace(request.DisplayName) == "" {
		return fmt.Errorf("display name is required")
	}

	return nil
}
