package services

import (
	"context"
	"testing"

	"github.com/igorracki/motorsports/backend/internal/auth"
	"github.com/igorracki/motorsports/backend/internal/database"
	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthService(t *testing.T) {
	databaseManager, err := database.NewManager(":memory:")
	require.NoError(t, err)
	defer databaseManager.Close()

	tokenManager, err := auth.NewTokenManager("test_secret_key_long_enough_for_jwt")
	require.NoError(t, err)

	userRepo := repository.NewUserRepository(databaseManager)
	authService := NewAuthService(userRepo, tokenManager)

	ctx := context.Background()

	t.Run("Register and Login - Success", func(tt *testing.T) {
		// Given
		registerRequest := models.RegisterUserRequest{
			Email:       "max@redbull.com",
			Password:    "password123",
			DisplayName: "Max Verstappen",
		}

		// When: Register
		user, profile, regToken, regExpires, err := authService.Register(ctx, registerRequest)
		require.NoError(tt, err)
		assert.NotNil(tt, user)
		assert.Equal(tt, registerRequest.Email, user.Email)
		assert.Equal(tt, registerRequest.DisplayName, profile.DisplayName)
		assert.NotEmpty(tt, regToken)
		assert.NotEmpty(tt, regExpires)

		// When: Login
		loginRequest := models.LoginRequest{
			Email:      registerRequest.Email,
			Password:   registerRequest.Password,
			RememberMe: true,
		}
		loggedInUser, loggedInProfile, token, expiresAt, err := authService.Login(ctx, loginRequest)

		// Then
		require.NoError(tt, err)
		assert.Equal(tt, user.ID, loggedInUser.ID)
		assert.Equal(tt, profile.DisplayName, loggedInProfile.DisplayName)
		assert.NotEmpty(tt, token)
		assert.NotEmpty(tt, expiresAt)
	})

	t.Run("Login - Invalid Password", func(tt *testing.T) {
		// Given
		registerRequest := models.RegisterUserRequest{
			Email:       "sergio@redbull.com",
			Password:    "password123",
			DisplayName: "Checo",
		}
		_, _, _, _, err := authService.Register(ctx, registerRequest)
		require.NoError(tt, err)

		// When
		loginRequest := models.LoginRequest{
			Email:    registerRequest.Email,
			Password: "wrong_password",
		}
		user, profile, token, _, err := authService.Login(ctx, loginRequest)

		// Then
		assert.Error(tt, err)
		assert.Nil(tt, user)
		assert.Nil(tt, profile)
		assert.Empty(tt, token)
		assert.Contains(tt, err.Error(), "invalid email or password")
	})
}
