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
	userService := NewUserService(userRepo, scoreRepo)

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
}
