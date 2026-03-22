package services

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/repository"
)

type LeaderboardService interface {
	GetLeaderboard(ctx context.Context, userID string, season int) ([]models.LeaderboardEntry, error)
}

type leaderboardService struct {
	friendRepo repository.FriendRepository
	userRepo   repository.UserRepository
	scoreRepo  repository.ScoreRepository
}

func NewLeaderboardService(friendRepo repository.FriendRepository, userRepo repository.UserRepository, scoreRepo repository.ScoreRepository) LeaderboardService {
	return &leaderboardService{
		friendRepo: friendRepo,
		userRepo:   userRepo,
		scoreRepo:  scoreRepo,
	}
}

func (service *leaderboardService) GetLeaderboard(ctx context.Context, userID string, season int) ([]models.LeaderboardEntry, error) {
	slog.InfoContext(ctx, "Entry: GetLeaderboard", "user_id", userID, "season", season)

	// Get friends
	friendIDs, err := service.friendRepo.GetFriendsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching friends: %w", err)
	}

	// Include self
	allUserIDs := make([]string, len(friendIDs)+1)
	copy(allUserIDs, friendIDs)
	allUserIDs[len(friendIDs)] = userID

	// Bulk fetch profiles
	profiles, err := service.userRepo.GetProfilesByUserIDs(ctx, allUserIDs)
	if err != nil {
		return nil, fmt.Errorf("fetching profiles: %w", err)
	}

	// Bulk fetch scores
	scores, err := service.scoreRepo.GetSeasonScoresByUserIDs(ctx, allUserIDs, season)
	if err != nil {
		return nil, fmt.Errorf("fetching scores: %w", err)
	}

	// Map scores for easy lookup
	userScores := make(map[string]int)
	for _, s := range scores {
		userScores[s.UserID] = s.Value
	}

	entries := make([]models.LeaderboardEntry, 0, len(profiles))
	for _, p := range profiles {
		entries = append(entries, models.LeaderboardEntry{
			UserID:      p.UserID,
			DisplayName: p.DisplayName,
			Score:       userScores[p.UserID], // Defaults to 0 if not found
		})
	}

	// Sort by score descending
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score != entries[j].Score {
			return entries[i].Score > entries[j].Score
		}
		return entries[i].DisplayName < entries[j].DisplayName
	})

	// Assign positions
	for i := range entries {
		entries[i].Position = i + 1
	}

	slog.InfoContext(ctx, "Exit: GetLeaderboard", "user_id", userID, "season", season, "count", len(entries))
	return entries, nil
}
