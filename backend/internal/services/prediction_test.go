package services

import (
	"context"
	"testing"
	"time"

	"github.com/igorracki/motorsports/backend/internal/database"
	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockF1Service struct {
	F1Service
	schedule []models.RaceWeekend
}

func (m *mockF1Service) GetScheduleByYear(ctx context.Context, year int) ([]models.RaceWeekend, error) {
	return m.schedule, nil
}

func TestPredictionService(t *testing.T) {
	databaseManager, err := database.NewManager(":memory:")
	require.NoError(t, err)
	defer databaseManager.Close()

	// Default future schedule for tests
	futureSchedule := []models.RaceWeekend{
		{
			Round: 1,
			Sessions: []models.Session{
				{Type: "Race", TimeUTCMS: time.Now().Add(24 * time.Hour).UnixMilli()},
			},
		},
		{
			Round: 2,
			Sessions: []models.Session{
				{Type: "Qualifying", TimeUTCMS: time.Now().Add(24 * time.Hour).UnixMilli()},
			},
		},
	}
	f1Mock := &mockF1Service{schedule: futureSchedule}

	predictionRepo := repository.NewPredictionRepository(databaseManager)
	userRepo := repository.NewUserRepository(databaseManager)
	scoringService := NewScoringService()
	predictionService := NewPredictionService(predictionRepo, f1Mock, scoringService)

	ctx := context.Background()

	// Helper to create a user for FK constraints
	createUser := func(id string) {
		user := &models.User{ID: id, Email: id + "@example.com"}
		profile := &models.Profile{UserID: id, DisplayName: "User " + id}
		err := userRepo.CreateUser(ctx, user, "hash", profile)
		require.NoError(t, err)
	}

	t.Run("Submit Prediction - Success", func(tt *testing.T) {
		// Given
		userID := "user-123"
		createUser(userID)
		prediction := &models.Prediction{
			UserID:      userID,
			Year:        2024,
			Round:       1,
			SessionType: "Race",
			Entries: []models.PredictionEntry{
				{Position: 1, DriverID: "VER"},
				{Position: 2, DriverID: "PER"},
				{Position: 3, DriverID: "ALO"},
			},
		}

		// When
		err := predictionService.SubmitPrediction(ctx, prediction)

		// Then
		require.NoError(tt, err)
		assert.NotEmpty(tt, prediction.ID)
		assert.NotZero(tt, prediction.CreatedAt)
		assert.NotZero(tt, prediction.UpdatedAt)

		// Verify it was saved
		saved, err := predictionRepo.GetPrediction(ctx, userID, 2024, 1, "Race")
		require.NoError(tt, err)
		assert.NotNil(tt, saved)
		assert.Equal(tt, 3, len(saved.Entries))
	})

	t.Run("Submit Prediction - Validation Errors", func(tt *testing.T) {
		userID := "user-val"
		createUser(userID)
		tests := []struct {
			name       string
			prediction *models.Prediction
			wantErr    string
		}{
			{
				name: "Too few entries",
				prediction: &models.Prediction{
					UserID: userID, Year: 2024, Round: 1, SessionType: "Race",
					Entries: []models.PredictionEntry{
						{Position: 1, DriverID: "VER"},
						{Position: 2, DriverID: "PER"},
					},
				},
				wantErr: "between 3 and 22 entries",
			},
			{
				name: "Duplicate position",
				prediction: &models.Prediction{
					UserID: userID, Year: 2024, Round: 1, SessionType: "Race",
					Entries: []models.PredictionEntry{
						{Position: 1, DriverID: "VER"},
						{Position: 1, DriverID: "PER"},
						{Position: 3, DriverID: "ALO"},
					},
				},
				wantErr: "duplicate position 1",
			},
			{
				name: "Duplicate driver",
				prediction: &models.Prediction{
					UserID: userID, Year: 2024, Round: 1, SessionType: "Race",
					Entries: []models.PredictionEntry{
						{Position: 1, DriverID: "VER"},
						{Position: 2, DriverID: "VER"},
						{Position: 3, DriverID: "ALO"},
					},
				},
				wantErr: "duplicate driver VER",
			},
			{
				name: "Invalid position (too high)",
				prediction: &models.Prediction{
					UserID: userID, Year: 2024, Round: 1, SessionType: "Race",
					Entries: []models.PredictionEntry{
						{Position: 1, DriverID: "VER"},
						{Position: 2, DriverID: "PER"},
						{Position: 23, DriverID: "ALO"},
					},
				},
				wantErr: "invalid position 23",
			},
			{
				name: "Empty driver ID",
				prediction: &models.Prediction{
					UserID: userID, Year: 2024, Round: 1, SessionType: "Race",
					Entries: []models.PredictionEntry{
						{Position: 1, DriverID: "VER"},
						{Position: 2, DriverID: ""},
						{Position: 3, DriverID: "ALO"},
					},
				},
				wantErr: "driver_id cannot be empty",
			},
		}

		for _, tc := range tests {
			tt.Run(tc.name, func(ttt *testing.T) {
				err := predictionService.SubmitPrediction(ctx, tc.prediction)
				assert.Error(ttt, err)
				assert.Contains(ttt, err.Error(), tc.wantErr)
			})
		}
	})

	t.Run("Submit Prediction - Update Existing", func(tt *testing.T) {
		// Given
		userID := "user-456"
		createUser(userID)
		initial := &models.Prediction{
			UserID:      userID,
			Year:        2024,
			Round:       2,
			SessionType: "Qualifying",
			Entries: []models.PredictionEntry{
				{Position: 1, DriverID: "VER"},
				{Position: 2, DriverID: "LEC"},
				{Position: 3, DriverID: "HAM"},
			},
		}
		err := predictionService.SubmitPrediction(ctx, initial)
		require.NoError(tt, err)
		initialID := initial.ID

		updated := &models.Prediction{
			UserID:      userID,
			Year:        2024,
			Round:       2,
			SessionType: "Qualifying",
			Entries: []models.PredictionEntry{
				{Position: 1, DriverID: "LEC"},
				{Position: 2, DriverID: "VER"},
				{Position: 3, DriverID: "HAM"},
			},
		}

		// When
		err = predictionService.SubmitPrediction(ctx, updated)

		// Then
		require.NoError(tt, err)
		assert.Equal(tt, initialID, updated.ID)
		assert.True(tt, initial.CreatedAt.Equal(updated.CreatedAt), "CreatedAt should be preserved")
		assert.True(tt, updated.UpdatedAt.After(initial.UpdatedAt), "UpdatedAt should advance")

		saved, err := predictionRepo.GetPrediction(ctx, userID, 2024, 2, "Qualifying")
		require.NoError(tt, err)
		assert.Equal(tt, "LEC", saved.Entries[0].DriverID)
	})

	t.Run("Get Predictions", func(tt *testing.T) {
		userID := "user-789"
		createUser(userID)
		p1 := &models.Prediction{
			UserID: userID, Year: 2024, Round: 1, SessionType: "Race",
			Entries: []models.PredictionEntry{{Position: 1, DriverID: "VER"}, {Position: 2, DriverID: "PER"}, {Position: 3, DriverID: "ALO"}},
		}
		require.NoError(tt, predictionService.SubmitPrediction(ctx, p1))

		// Get all
		all, err := predictionService.GetUserPredictions(ctx, userID)
		require.NoError(tt, err)
		assert.Len(tt, all, 1)

		// Get round
		round, err := predictionService.GetRoundPredictions(ctx, userID, 2024, 1)
		require.NoError(tt, err)
		assert.Len(tt, round, 1)
		assert.Equal(tt, "Race", round[0].SessionType)
	})
}
