package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/igorracki/f1/backend/internal/database"
	"github.com/igorracki/f1/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPredictionRepository(t *testing.T) {
	databaseManager, err := database.NewManager(":memory:")
	require.NoError(t, err)
	defer databaseManager.Close()

	userRepo := NewUserRepository(databaseManager.DB())
	predictionRepo := NewPredictionRepository(databaseManager.DB())

	// Create a user first for FK constraints
	userID := uuid.New().String()
	err = userRepo.CreateUser(&models.User{
		ID: userID, Username: "max33", Email: "max@redbull.com", CreatedAt: time.Now(),
	}, &models.Profile{UserID: userID, DisplayName: "Max"})
	require.NoError(t, err)

	t.Run("Save and Get Prediction", func(tt *testing.T) {
		// Given
		predictionID := uuid.New().String()
		prediction := &models.Prediction{
			ID:          predictionID,
			UserID:      userID,
			Year:        2024,
			Round:       1,
			SessionType: "race",
			CreatedAt:   time.Now().UTC().Truncate(time.Second),
			UpdatedAt:   time.Now().UTC().Truncate(time.Second),
			Entries: []models.PredictionEntry{
				{Position: 1, DriverID: "verstappen"},
				{Position: 2, DriverID: "perez"},
			},
		}

		// When
		err := predictionRepo.SavePrediction(prediction)
		require.NoError(tt, err)

		fetched, err := predictionRepo.GetPrediction(userID, 2024, 1, "race")
		require.NoError(tt, err)

		// Then
		assert.NotNil(tt, fetched)
		assert.Equal(tt, 2, len(fetched.Entries))
		assert.Equal(tt, "verstappen", fetched.Entries[0].DriverID)
		assert.Equal(tt, 1, fetched.Entries[0].Position)
	})

	t.Run("Update Prediction (Upsert)", func(tt *testing.T) {
		// Given: Change P2 from Perez to Hamilton
		updatedPrediction := &models.Prediction{
			ID:          uuid.New().String(), // New ID but same User/Year/Round/Session
			UserID:      userID,
			Year:        2024,
			Round:       1,
			SessionType: "race",
			CreatedAt:   time.Now().UTC().Truncate(time.Second),
			UpdatedAt:   time.Now().UTC().Truncate(time.Second),
			Entries: []models.PredictionEntry{
				{Position: 1, DriverID: "verstappen"},
				{Position: 2, DriverID: "hamilton"},
			},
		}

		// When
		err := predictionRepo.SavePrediction(updatedPrediction)
		require.NoError(tt, err)

		fetched, err := predictionRepo.GetPrediction(userID, 2024, 1, "race")
		require.NoError(tt, err)

		// Then
		assert.Equal(tt, 2, len(fetched.Entries))
		assert.Equal(tt, "hamilton", fetched.Entries[1].DriverID)
	})
}
