package services

import (
	"context"
	"fmt"
	"sort"

	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/repository"
)

type LeaderboardService interface {
	GetLeaderboard(ctx context.Context, userID string, season int) ([]models.LeaderboardEntry, error)
}

type leaderboardService struct {
	friendRepo     repository.FriendRepository
	userRepo       repository.UserRepository
	predictionRepo repository.PredictionRepository
}

func NewLeaderboardService(friendRepo repository.FriendRepository, userRepo repository.UserRepository, predictionRepo repository.PredictionRepository) LeaderboardService {
	return &leaderboardService{
		friendRepo:     friendRepo,
		userRepo:       userRepo,
		predictionRepo: predictionRepo,
	}
}

func (service *leaderboardService) GetLeaderboard(ctx context.Context, userID string, season int) ([]models.LeaderboardEntry, error) {
	friendIDs, err := service.friendRepo.GetFriendsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching friends: %w", err)
	}

	allUserIDs := make([]string, len(friendIDs)+1)
	copy(allUserIDs, friendIDs)
	allUserIDs[len(friendIDs)] = userID

	profiles, err := service.userRepo.GetProfilesByUserIDs(ctx, allUserIDs)
	if err != nil {
		return nil, fmt.Errorf("fetching profiles: %w", err)
	}

	userScores, err := service.predictionRepo.GetSeasonScoresByUserIDs(ctx, allUserIDs, season)
	if err != nil {
		return nil, fmt.Errorf("fetching aggregated scores from predictions: %w", err)
	}

	entries := make([]models.LeaderboardEntry, 0, len(profiles))
	for _, profile := range profiles {
		entries = append(entries, models.LeaderboardEntry{
			UserID:      profile.UserID,
			DisplayName: profile.DisplayName,
			Score:       userScores[profile.UserID],
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score != entries[j].Score {
			return entries[i].Score > entries[j].Score
		}
		return entries[i].DisplayName < entries[j].DisplayName
	})

	for i := range entries {
		entries[i].Position = i + 1
	}

	return entries, nil
}
