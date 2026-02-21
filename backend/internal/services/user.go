package services

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/repository"
	"github.com/microcosm-cc/bluemonday"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
var sanitizer = bluemonday.StrictPolicy()

type UserService interface {
	RegisterUser(ctx context.Context, request models.RegisterUserRequest) (*models.User, error)
	GetFullProfile(ctx context.Context, userID string) (*models.UserProfileResponse, error)
}

type userService struct {
	userRepository  repository.UserRepository
	scoreRepository repository.ScoreRepository
}

func NewUserService(userRepo repository.UserRepository, scoreRepo repository.ScoreRepository) UserService {
	return &userService{
		userRepository:  userRepo,
		scoreRepository: scoreRepo,
	}
}

func (service *userService) RegisterUser(ctx context.Context, request models.RegisterUserRequest) (*models.User, error) {
	log.Printf("INFO: Attempting to register user [username: %s, email: %s]", request.Username, request.Email)

	if err := service.validateRegistrationRequest(request); err != nil {
		log.Printf("WARN: Registration validation failed [username: %s, error: %v]", request.Username, err)
		return nil, err
	}

	userID := uuid.New().String()
	user := &models.User{
		ID:        userID,
		Username:  request.Username,
		Email:     request.Email,
		CreatedAt: time.Now().UTC(),
	}

	displayName := request.DisplayName
	if strings.TrimSpace(displayName) == "" {
		displayName = request.Username
	}
	displayName = sanitizer.Sanitize(displayName)

	profile := &models.Profile{
		UserID:      userID,
		DisplayName: displayName,
	}

	if err := service.userRepository.CreateUser(ctx, user, profile); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			if strings.Contains(err.Error(), "users.username") {
				return nil, fmt.Errorf("username '%s' is already taken", request.Username)
			}
			if strings.Contains(err.Error(), "users.email") {
				return nil, fmt.Errorf("email '%s' is already registered", request.Email)
			}
		}
		return nil, fmt.Errorf("failed to create user in database: %w", err)
	}

	log.Printf("INFO: Successfully registered user [id: %s, username: %s]", userID, user.Username)
	return user, nil
}

func (service *userService) GetFullProfile(ctx context.Context, userID string) (*models.UserProfileResponse, error) {
	log.Printf("INFO: Fetching full profile for user [id: %s]", userID)

	user, err := service.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching user %s: %w", userID, err)
	}
	if user == nil {
		log.Printf("WARN: User not found while fetching profile [id: %s]", userID)
		return nil, fmt.Errorf("user not found")
	}

	profile, err := service.userRepository.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching profile for user %s: %w", userID, err)
	}
	if profile == nil {
		log.Printf("WARN: Profile not found for user [id: %s]", userID)
		return nil, fmt.Errorf("profile not found")
	}

	scores, err := service.scoreRepository.GetUserScores(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching scores for user %s: %w", userID, err)
	}

	if scores == nil {
		scores = []models.UserScore{}
	}

	response := &models.UserProfileResponse{
		User:    *user,
		Profile: *profile,
		Scores:  scores,
	}

	log.Printf("INFO: Successfully fetched full profile [id: %s, scores_count: %d]", userID, len(scores))
	return response, nil
}

func (service *userService) validateRegistrationRequest(request models.RegisterUserRequest) error {
	if !usernameRegex.MatchString(request.Username) {
		return fmt.Errorf("invalid username: must be 3-20 characters and only contain letters, numbers, and underscores")
	}

	if _, err := mail.ParseAddress(request.Email); err != nil {
		return fmt.Errorf("invalid email address format")
	}

	return nil
}
