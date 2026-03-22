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

func TestUserRepository(t *testing.T) {
	databaseManager, err := database.NewManager(":memory:")
	require.NoError(t, err)
	defer databaseManager.Close()

	userRepo := NewUserRepository(databaseManager.DB())
	ctx := context.Background()

	t.Run("Create and Get User", func(tt *testing.T) {
		// Given
		userID := uuid.New().String()
		user := &models.User{
			ID:        userID,
			Email:     "lewis@mercedes.com",
			CreatedAt: time.Now().UTC().Truncate(time.Second),
		}
		profile := &models.Profile{
			UserID:      userID,
			DisplayName: "Sir Lewis",
		}

		// When
		err := userRepo.CreateUser(ctx, user, "hash", profile)
		require.NoError(tt, err)

		fetchedUser, err := userRepo.GetUserByID(ctx, userID)
		require.NoError(tt, err)

		// Then
		assert.NotNil(tt, fetchedUser)
		assert.Equal(tt, user.Email, fetchedUser.Email)
		assert.True(tt, user.CreatedAt.Equal(fetchedUser.CreatedAt))
	})

	t.Run("Get Non-Existent User", func(tt *testing.T) {
		// When
		fetchedUser, err := userRepo.GetUserByID(ctx, "non-existent")

		// Then
		require.NoError(tt, err)
		assert.Nil(tt, fetchedUser)
	})

	t.Run("Duplicate Constraints", func(tt *testing.T) {
		// Given
		userID1 := uuid.New().String()
		user1 := &models.User{
			ID:        userID1,
			Email:     "charles@ferrari.com",
			CreatedAt: time.Now().UTC().Truncate(time.Second),
		}
		profile1 := &models.Profile{UserID: userID1, DisplayName: "Charles"}
		err := userRepo.CreateUser(ctx, user1, "hash", profile1)
		require.NoError(tt, err)

		// When: Duplicate email
		userID3 := uuid.New().String()
		user3 := &models.User{
			ID:        userID3,
			Email:     "charles@ferrari.com", // Duplicate
			CreatedAt: time.Now().UTC().Truncate(time.Second),
		}
		profile3 := &models.Profile{UserID: userID3, DisplayName: "Carlos"}
		err = userRepo.CreateUser(ctx, user3, "hash", profile3)

		// Then
		assert.Error(tt, err)
		assert.ErrorIs(tt, err, ErrDuplicateEmail)
	})
}
