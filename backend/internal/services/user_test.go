package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/igorracki/f1/backend/internal/database"
	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService(t *testing.T) {
	databaseManager, err := database.NewManager(":memory:")
	require.NoError(t, err)
	defer databaseManager.Close()

	userRepo := repository.NewUserRepository(databaseManager.DB())
	scoreRepo := repository.NewScoreRepository(databaseManager.DB())
	predictionRepo := repository.NewPredictionRepository(databaseManager.DB())

	f1Mock := &mockF1Service{} // Need a basic mock or real service if tests call it
	scoringService := NewScoringService()
	predictionService := NewPredictionService(predictionRepo, f1Mock, scoringService)

	userService := NewUserService(userRepo, scoreRepo, predictionService)

	ctx := context.Background()

	t.Run("Get Full Profile - Success", func(tt *testing.T) {
		// Given
		userID := uuid.New().String()
		user := &models.User{
			ID:        userID,
			Email:     "oscar@mclaren.com",
			CreatedAt: time.Now(),
		}
		profile := &models.Profile{
			UserID:      userID,
			DisplayName: "Oscar Piastri",
		}
		err := userRepo.CreateUser(ctx, user, "hash", profile)
		require.NoError(tt, err)

		// When
		fetchedProfile, err := userService.GetFullProfile(ctx, userID)

		// Then
		require.NoError(tt, err)
		assert.NotNil(tt, fetchedProfile)
		assert.Equal(tt, userID, fetchedProfile.User.ID)
		assert.Equal(tt, "Oscar Piastri", fetchedProfile.Profile.DisplayName)
		assert.Empty(tt, fetchedProfile.Scores)
	})

	t.Run("Get Full Profile - With Aggregated Scores", func(tt *testing.T) {
		// Given
		userID := uuid.New().String()
		user := &models.User{ID: userID, Email: "lando@mclaren.com", CreatedAt: time.Now()}
		profile := &models.Profile{UserID: userID, DisplayName: "Lando"}
		require.NoError(tt, userRepo.CreateUser(ctx, user, "hash", profile))

		// Add predictions with scores
		score1, score2 := 25, 18
		p1 := &models.Prediction{
			ID: uuid.New().String(), UserID: userID, Year: 2024, Round: 1, SessionType: "Race", Score: &score1,
			Entries: []models.PredictionEntry{{Position: 1, DriverID: "VER"}},
		}
		p2 := &models.Prediction{
			ID: uuid.New().String(), UserID: userID, Year: 2024, Round: 2, SessionType: "Race", Score: &score2,
			Entries: []models.PredictionEntry{{Position: 1, DriverID: "LEC"}},
		}
		require.NoError(tt, predictionRepo.SavePrediction(ctx, p1))
		require.NoError(tt, predictionRepo.SavePrediction(ctx, p2))

		// When
		fetchedScores, err := userService.GetSeasonScores(ctx, userID)

		// Then
		require.NoError(tt, err)
		assert.Len(tt, fetchedScores, 1)
		assert.Equal(tt, "season_total", fetchedScores[0].ScoreType)
		assert.Equal(tt, 2024, *fetchedScores[0].Season)
		assert.Equal(tt, 43, fetchedScores[0].Value)
	})
}
