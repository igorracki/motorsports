package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/igorracki/motorsports/backend/internal/auth"
	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/repository"
)

type AuthService interface {
	Register(ctx context.Context, request models.RegisterUserRequest) (*models.User, *models.Profile, string, time.Time, error)
	Login(ctx context.Context, request models.LoginRequest) (*models.User, *models.Profile, string, time.Time, error)
}

type authService struct {
	userRepository repository.UserRepository
	tokenManager   auth.TokenManager
}

func NewAuthService(userRepo repository.UserRepository, tokenManager auth.TokenManager) AuthService {
	return &authService{
		userRepository: userRepo,
		tokenManager:   tokenManager,
	}
}

func (service *authService) Register(ctx context.Context, request models.RegisterUserRequest) (*models.User, *models.Profile, string, time.Time, error) {
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

	displayName := request.DisplayName
	if strings.TrimSpace(displayName) == "" {
		displayName = strings.Split(request.Email, "@")[0]
	}

	profile := &models.Profile{
		UserID:      userID,
		DisplayName: displayName,
	}

	if err := service.userRepository.CreateUser(ctx, user, passwordHash, profile); err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return nil, nil, "", time.Time{}, fmt.Errorf("%w: email is already registered", models.ErrConflict)
		}
		return nil, nil, "", time.Time{}, fmt.Errorf("failed to create user: %w", err)
	}

	token, expiresAt, err := service.tokenManager.GenerateToken(userID, 24*time.Hour)
	if err != nil {
		return user, profile, "", time.Time{}, fmt.Errorf("generating token: %w", err)
	}

	slog.InfoContext(ctx, "User registered successfully", "user_id", userID)
	return user, profile, token, expiresAt, nil
}

func (service *authService) Login(ctx context.Context, request models.LoginRequest) (*models.User, *models.Profile, string, time.Time, error) {
	user, passwordHash, err := service.userRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, nil, "", time.Time{}, fmt.Errorf("fetching user: %w", err)
	}

	if user == nil {
		return nil, nil, "", time.Time{}, fmt.Errorf("%w: invalid email or password", models.ErrUnauthorized)
	}

	match, err := auth.VerifyPassword(request.Password, passwordHash)
	if err != nil {
		return nil, nil, "", time.Time{}, fmt.Errorf("verifying password: %w", err)
	}

	if !match {
		return nil, nil, "", time.Time{}, fmt.Errorf("%w: invalid email or password", models.ErrUnauthorized)
	}

	profile, err := service.userRepository.GetProfileByUserID(ctx, user.ID)
	if err != nil {
		return nil, nil, "", time.Time{}, fmt.Errorf("fetching profile: %w", err)
	}

	duration := 24 * time.Hour
	if request.RememberMe {
		duration = 30 * 24 * time.Hour
	}

	token, expiresAt, err := service.tokenManager.GenerateToken(user.ID, duration)
	if err != nil {
		return nil, nil, "", time.Time{}, fmt.Errorf("generating token: %w", err)
	}

	slog.InfoContext(ctx, "User logged in successfully", "user_id", user.ID)
	return user, profile, token, expiresAt, nil
}
