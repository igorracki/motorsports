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
	GetSeasonScores(ctx context.Context, userID string) ([]models.UserScore, error)
}

type userService struct {
	userRepository    repository.UserRepository
	scoreRepository   repository.ScoreRepository
	predictionService PredictionService
}

func NewUserService(userRepo repository.UserRepository, scoreRepo repository.ScoreRepository, predictionService PredictionService) UserService {
	return &userService{
		userRepository:    userRepo,
		scoreRepository:   scoreRepo,
		predictionService: predictionService,
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

	// Fetch explicit scores from score repository
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

	log.Printf("INFO: Successfully fetched basic profile [id: %s, scores_count: %d]", userID, len(scores))
	return response, nil
}

func (service *userService) GetSeasonScores(ctx context.Context, userID string) ([]models.UserScore, error) {
	log.Printf("INFO: Calculating season scores for user [id: %s]", userID)

	predictions, err := service.predictionService.GetUserPredictions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching predictions for user %s: %w", userID, err)
	}

	seasonTotals := make(map[int]int)
	for _, p := range predictions {
		if p.Score != nil {
			seasonTotals[p.Year] += *p.Score
		}
	}

	scores := make([]models.UserScore, 0, len(seasonTotals))
	for year, total := range seasonTotals {
		yearCopy := year
		scores = append(scores, models.UserScore{
			UserID:    userID,
			ScoreType: "season_total",
			Season:    &yearCopy,
			Value:     total,
		})
	}

	log.Printf("INFO: Calculated scores for %d seasons for user [id: %s]", len(scores), userID)
	return scores, nil
}
