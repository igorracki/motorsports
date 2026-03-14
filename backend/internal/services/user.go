package services

import (
	"context"
	"fmt"
	"log"

	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/repository"
)

type UserService interface {
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
