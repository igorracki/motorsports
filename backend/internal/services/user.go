package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/repository"
)

type UserService interface {
	GetFullProfile(ctx context.Context, userID string) (*models.UserProfileResponse, error)
	GetSeasonScores(ctx context.Context, userID string) ([]models.UserScore, error)
}

type userService struct {
	userRepository    repository.UserRepository
	predictionService PredictionService
}

func NewUserService(userRepo repository.UserRepository, predictionService PredictionService) UserService {
	return &userService{
		userRepository:    userRepo,
		predictionService: predictionService,
	}
}

func (service *userService) GetFullProfile(ctx context.Context, userID string) (*models.UserProfileResponse, error) {
	slog.InfoContext(ctx, "Entry: GetFullProfile", "user_id", userID)

	user, err := service.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching user %s: %w", userID, err)
	}
	if user == nil {
		slog.WarnContext(ctx, "User not found while fetching profile", "user_id", userID)
		return nil, fmt.Errorf("user not found")
	}

	profile, err := service.userRepository.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching profile for user %s: %w", userID, err)
	}
	if profile == nil {
		slog.WarnContext(ctx, "Profile not found for user", "user_id", userID)
		return nil, fmt.Errorf("profile not found")
	}

	scores, err := service.GetSeasonScores(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching aggregated scores for user %s: %w", userID, err)
	}

	response := &models.UserProfileResponse{
		User:    *user,
		Profile: *profile,
		Scores:  scores,
	}

	slog.InfoContext(ctx, "Exit: GetFullProfile", "user_id", userID, "scores_count", len(scores))
	return response, nil
}

func (service *userService) GetSeasonScores(ctx context.Context, userID string) ([]models.UserScore, error) {
	slog.InfoContext(ctx, "Entry: GetSeasonScores", "user_id", userID)

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

	slog.InfoContext(ctx, "Exit: GetSeasonScores", "user_id", userID, "seasons_count", len(scores))
	return scores, nil
}
