package services

import (
	"context"
	"testing"

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

	t.Run("Register User - Success", func(tt *testing.T) {
		// Given
		request := models.RegisterUserRequest{
			Username:    "valtteri77",
			Email:       "valtteri@sauber.com",
			DisplayName: "Valtteri Bottas",
		}

		// When
		user, err := userService.RegisterUser(ctx, request)

		// Then
		require.NoError(tt, err)
		assert.NotNil(tt, user)
		assert.Equal(tt, request.Username, user.Username)
		assert.Equal(tt, request.Email, user.Email)
		assert.NotEmpty(tt, user.ID)
	})

	t.Run("Register User - Invalid Username", func(tt *testing.T) {
		// Given
		request := models.RegisterUserRequest{
			Username:    "v", // Too short
			Email:       "v@b.com",
			DisplayName: "V",
		}

		// When
		user, err := userService.RegisterUser(ctx, request)

		// Then
		assert.Error(tt, err)
		assert.Nil(tt, user)
		assert.Contains(tt, err.Error(), "invalid username")
	})

	t.Run("Register User - Duplicate Username", func(tt *testing.T) {
		// Given
		request1 := models.RegisterUserRequest{
			Username:    "george63",
			Email:       "george@mercedes.com",
			DisplayName: "George Russell",
		}
		_, err := userService.RegisterUser(ctx, request1)
		require.NoError(tt, err)

		request2 := models.RegisterUserRequest{
			Username:    "george63", // Duplicate
			Email:       "other@mercedes.com",
			DisplayName: "Other George",
		}

		// When
		user, err := userService.RegisterUser(ctx, request2)

		// Then
		assert.Error(tt, err)
		assert.Nil(tt, user)
		assert.Contains(tt, err.Error(), "already taken")
	})

	t.Run("Get Full Profile - Success", func(tt *testing.T) {
		// Given
		request := models.RegisterUserRequest{
			Username:    "oscar81",
			Email:       "oscar@mclaren.com",
			DisplayName: "Oscar Piastri",
		}
		user, err := userService.RegisterUser(ctx, request)
		require.NoError(tt, err)

		// When
		profile, err := userService.GetFullProfile(ctx, user.ID)

		// Then
		require.NoError(tt, err)
		assert.NotNil(tt, profile)
		assert.Equal(tt, user.ID, profile.User.ID)
		assert.Equal(tt, "Oscar Piastri", profile.Profile.DisplayName)
		assert.Empty(tt, profile.Scores)
	})

	t.Run("Register User - DisplayName Fallback", func(tt *testing.T) {
		// Given
		request := models.RegisterUserRequest{
			Username: "nico27",
			Email:    "nico@haas.com",
			// DisplayName is empty
		}

		// When
		user, err := userService.RegisterUser(ctx, request)
		require.NoError(tt, err)

		profile, err := userService.GetFullProfile(ctx, user.ID)
		require.NoError(tt, err)

		// Then
		assert.Equal(tt, "nico27", profile.Profile.DisplayName)
	})

	t.Run("Register User - HTML Sanitization", func(tt *testing.T) {
		// Given
		request := models.RegisterUserRequest{
			Username:    "hack_erman",
			Email:       "hacker@evil.com",
			DisplayName: "<script>alert('XSS')</script>Hacker",
		}

		// When
		user, err := userService.RegisterUser(ctx, request)
		require.NoError(tt, err)

		profile, err := userService.GetFullProfile(ctx, user.ID)
		require.NoError(tt, err)

		// Then
		assert.Equal(tt, "Hacker", profile.Profile.DisplayName)
	})
}
