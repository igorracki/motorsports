package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/igorracki/motorsports/backend/internal/database"
	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPredictionRepository(t *testing.T) {
	databaseManager, err := database.NewManager(":memory:")
	require.NoError(t, err)
	defer databaseManager.Close()

	userRepo := NewUserRepository(databaseManager.DB())
	predictionRepo := NewPredictionRepository(databaseManager.DB())
	ctx := context.Background()

	// Create a user first for FK constraints
	userID := uuid.New().String()
	err = userRepo.CreateUser(ctx, &models.User{
		ID: userID, Email: "max@redbull.com", CreatedAt: time.Now(),
	}, "hash", &models.Profile{UserID: userID, DisplayName: "Max"})
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
		err := predictionRepo.SavePrediction(ctx, prediction)
		require.NoError(tt, err)

		fetched, err := predictionRepo.GetPrediction(ctx, userID, 2024, 1, "race")
		require.NoError(tt, err)

		// Then
		assert.NotNil(tt, fetched)
		assert.Equal(tt, 2, len(fetched.Entries))
		assert.Equal(tt, "verstappen", fetched.Entries[0].DriverID)
		assert.Equal(tt, 1, fetched.Entries[0].Position)
	})

	t.Run("Update Prediction (Upsert)", func(tt *testing.T) {
		// Given: First submission
		initialPrediction := &models.Prediction{
			ID:          uuid.New().String(),
			UserID:      userID,
			Year:        2025,
			Round:       1,
			SessionType: "race",
			CreatedAt:   time.Now().UTC().Truncate(time.Second),
			UpdatedAt:   time.Now().UTC().Truncate(time.Second),
			Entries: []models.PredictionEntry{
				{Position: 1, DriverID: "verstappen"},
				{Position: 2, DriverID: "perez"},
			},
		}
		err := predictionRepo.SavePrediction(ctx, initialPrediction)
		require.NoError(tt, err)
		initialID := initialPrediction.ID
		initialCreatedAt := initialPrediction.CreatedAt

		// When: Change P2 from Perez to Hamilton
		updatedPrediction := &models.Prediction{
			ID:          uuid.New().String(), // This should be overwritten by the original ID
			UserID:      userID,
			Year:        2025,
			Round:       1,
			SessionType: "race",
			CreatedAt:   time.Now().UTC().Truncate(time.Second),
			UpdatedAt:   time.Now().UTC().Truncate(time.Second),
			Entries: []models.PredictionEntry{
				{Position: 1, DriverID: "verstappen"},
				{Position: 2, DriverID: "hamilton"},
			},
		}

		err = predictionRepo.SavePrediction(ctx, updatedPrediction)
		require.NoError(tt, err)

		fetched, err := predictionRepo.GetPrediction(ctx, userID, 2025, 1, "race")
		require.NoError(tt, err)

		// Then
		assert.Equal(tt, initialID, fetched.ID, "ID should be preserved")
		assert.Equal(tt, initialID, updatedPrediction.ID, "Input model ID should be updated to original ID")
		assert.True(tt, initialCreatedAt.Equal(fetched.CreatedAt), "CreatedAt should be preserved")
		assert.Equal(tt, 2, len(fetched.Entries))
		assert.Equal(tt, "hamilton", fetched.Entries[1].DriverID)
	})

	t.Run("Get User and Round Predictions", func(tt *testing.T) {
		// Given: Two predictions for different rounds
		p1 := &models.Prediction{
			ID: uuid.New().String(), UserID: userID, Year: 2024, Round: 1, SessionType: "Qualifying",
			CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
			Entries: []models.PredictionEntry{{Position: 1, DriverID: "VER"}},
		}
		p2 := &models.Prediction{
			ID: uuid.New().String(), UserID: userID, Year: 2024, Round: 2, SessionType: "Race",
			CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
			Entries: []models.PredictionEntry{{Position: 1, DriverID: "LEC"}},
		}
		require.NoError(tt, predictionRepo.SavePrediction(ctx, p1))
		require.NoError(tt, predictionRepo.SavePrediction(ctx, p2))

		// When: Fetching all predictions
		all, err := predictionRepo.GetUserPredictions(ctx, userID)
		require.NoError(tt, err)

		// Then: Should have at least 2 predictions from this subtest
		assert.GreaterOrEqual(tt, len(all), 2)

		// When: Fetching specific round
		round1, err := predictionRepo.GetRoundPredictions(ctx, userID, 2024, 1)
		require.NoError(tt, err)

		// Then: Should find the qualifying prediction
		assert.Len(tt, round1, 2) // p1 and the one from the first test case
		foundP1 := false
		for _, p := range round1 {
			if p.ID == p1.ID {
				foundP1 = true
				assert.Equal(tt, 1, len(p.Entries))
				assert.Equal(tt, "VER", p.Entries[0].DriverID)
			}
		}
		assert.True(tt, foundP1)
	})

	t.Run("Empty Results return empty slice not nil", func(tt *testing.T) {
		// Given: A new user with no predictions
		newUserID := uuid.New().String()
		err = userRepo.CreateUser(ctx, &models.User{
			ID: newUserID, Email: "empty@example.com", CreatedAt: time.Now(),
		}, "hash", &models.Profile{UserID: newUserID, DisplayName: "Empty"})
		require.NoError(tt, err)

		// When
		all, err := predictionRepo.GetUserPredictions(ctx, newUserID)
		require.NoError(tt, err)
		round, err := predictionRepo.GetRoundPredictions(ctx, newUserID, 2024, 1)
		require.NoError(tt, err)

		// Then
		assert.NotNil(tt, all)
		assert.Equal(tt, 0, len(all))
		assert.NotNil(tt, round)
		assert.Equal(tt, 0, len(round))
	})
}
