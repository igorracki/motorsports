package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/repository"
)

type PredictionService interface {
	SubmitPrediction(ctx context.Context, prediction *models.Prediction) error
	GetUserPredictions(ctx context.Context, userID string) ([]models.Prediction, error)
	GetRoundPredictions(ctx context.Context, userID string, year, round int) ([]models.Prediction, error)
}

type predictionService struct {
	predictionRepository repository.PredictionRepository
}

func NewPredictionService(repo repository.PredictionRepository) PredictionService {
	return &predictionService{
		predictionRepository: repo,
	}
}

func (service *predictionService) SubmitPrediction(ctx context.Context, prediction *models.Prediction) error {
	slog.InfoContext(ctx, "Entry: SubmitPrediction",
		"user_id", prediction.UserID,
		"year", prediction.Year,
		"round", prediction.Round,
		"session", prediction.SessionType,
		"entry_count", len(prediction.Entries))

	if err := service.validatePrediction(prediction); err != nil {
		slog.WarnContext(ctx, "Prediction validation failed", "error", err)
		return err
	}

	now := time.Now().UTC()
	prediction.ID = uuid.New().String()
	prediction.CreatedAt = now
	prediction.UpdatedAt = now

	if err := service.predictionRepository.SavePrediction(ctx, prediction); err != nil {
		return fmt.Errorf("saving prediction to database: %w", err)
	}

	slog.InfoContext(ctx, "Exit: SubmitPrediction",
		"user_id", prediction.UserID,
		"prediction_id", prediction.ID)
	return nil
}

func (service *predictionService) GetUserPredictions(ctx context.Context, userID string) ([]models.Prediction, error) {
	slog.InfoContext(ctx, "Entry: GetUserPredictions", "user_id", userID)

	predictions, err := service.predictionRepository.GetUserPredictions(ctx, userID)
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "Exit: GetUserPredictions", "user_id", userID, "count", len(predictions))
	return predictions, nil
}

func (service *predictionService) GetRoundPredictions(ctx context.Context, userID string, year, round int) ([]models.Prediction, error) {
	slog.InfoContext(ctx, "Entry: GetRoundPredictions", "user_id", userID, "year", year, "round", round)

	predictions, err := service.predictionRepository.GetRoundPredictions(ctx, userID, year, round)
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "Exit: GetRoundPredictions", "user_id", userID, "year", year, "round", round, "count", len(predictions))
	return predictions, nil
}

func (service *predictionService) validatePrediction(prediction *models.Prediction) error {
	if len(prediction.Entries) < 3 || len(prediction.Entries) > 22 {
		return fmt.Errorf("prediction must have between 3 and 22 entries, got %d", len(prediction.Entries))
	}

	if prediction.Year < 1950 || prediction.Year > 2100 {
		return fmt.Errorf("invalid year: %d", prediction.Year)
	}

	if prediction.Round < 1 || prediction.Round > 50 {
		return fmt.Errorf("invalid round: %d", prediction.Round)
	}

	positions := make(map[int]bool)
	drivers := make(map[string]bool)

	for _, entry := range prediction.Entries {
		if entry.Position < 1 || entry.Position > 22 {
			return fmt.Errorf("invalid position %d: must be between 1 and 22", entry.Position)
		}
		if positions[entry.Position] {
			return fmt.Errorf("duplicate position %d in prediction", entry.Position)
		}
		positions[entry.Position] = true

		if entry.DriverID == "" {
			return fmt.Errorf("driver_id cannot be empty")
		}
		if drivers[entry.DriverID] {
			return fmt.Errorf("duplicate driver %s in prediction", entry.DriverID)
		}
		drivers[entry.DriverID] = true
	}

	return nil
}
