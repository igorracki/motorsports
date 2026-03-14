package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/igorracki/f1/backend/internal/database"
	"github.com/igorracki/f1/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScoreRepository(t *testing.T) {
	databaseManager, err := database.NewManager(":memory:")
	require.NoError(t, err)
	defer databaseManager.Close()

	userRepo := NewUserRepository(databaseManager.DB())
	scoreRepo := NewScoreRepository(databaseManager.DB())
	ctx := context.Background()

	userID := uuid.New().String()
	err = userRepo.CreateUser(ctx, &models.User{
		ID: userID, Email: "lando@mclaren.com", CreatedAt: time.Now(),
	}, "hash", &models.Profile{UserID: userID, DisplayName: "Lando"})
	require.NoError(t, err)

	t.Run("Update and Get Scores", func(tt *testing.T) {
		// Given
		season := 2024
		score := &models.UserScore{
			UserID:    userID,
			ScoreType: "points",
			Season:    &season,
			Value:     150,
			UpdatedAt: time.Now().UTC().Truncate(time.Second),
		}

		// When
		err := scoreRepo.UpdateScore(ctx, score)
		require.NoError(tt, err)

		// Then
		scores, err := scoreRepo.GetUserScores(ctx, userID)
		require.NoError(tt, err)
		assert.Equal(tt, 1, len(scores))
		assert.Equal(tt, 150, scores[0].Value)

		// When: Update value
		score.Value = 200
		err = scoreRepo.UpdateScore(ctx, score)
		require.NoError(tt, err)

		// Then
		scores, err = scoreRepo.GetUserScores(ctx, userID)
		require.NoError(tt, err)
		assert.Equal(tt, 1, len(scores))
		assert.Equal(tt, 200, scores[0].Value)
	})
}
