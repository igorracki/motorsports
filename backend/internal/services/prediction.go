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
	f1Service            F1Service
	scoringService       ScoringService
}

func NewPredictionService(repo repository.PredictionRepository, f1 F1Service, scoring ScoringService) PredictionService {
	return &predictionService{
		predictionRepository: repo,
		f1Service:            f1,
		scoringService:       scoring,
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
		return fmt.Errorf("validating prediction: %w", err)
	}

	sessionTimeMS, err := service.getSessionStartTimeMS(ctx, prediction.Year, prediction.Round, prediction.SessionType)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get session time for deadline check", "error", err)
		return fmt.Errorf("getting session start time: %w", err)
	}

	if time.Now().UTC().UnixMilli() >= sessionTimeMS {
		return fmt.Errorf("prediction period has closed for this session")
	}

	now := time.Now().UTC()
	revalidateUntil := time.UnixMilli(sessionTimeMS + (48 * 3600 * 1000)).UTC()
	prediction.ID = uuid.New().String()
	prediction.RevalidateUntil = &revalidateUntil
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
		return nil, fmt.Errorf("fetching user predictions from repository: %w", err)
	}

	for i := range predictions {
		service.processPredictionScore(ctx, &predictions[i])
	}

	slog.InfoContext(ctx, "Exit: GetUserPredictions", "user_id", userID, "count", len(predictions))
	return predictions, nil
}

func (service *predictionService) GetRoundPredictions(ctx context.Context, userID string, year, round int) ([]models.Prediction, error) {
	slog.InfoContext(ctx, "Entry: GetRoundPredictions", "user_id", userID, "year", year, "round", round)

	predictions, err := service.predictionRepository.GetRoundPredictions(ctx, userID, year, round)
	if err != nil {
		return nil, fmt.Errorf("fetching round predictions from repository: %w", err)
	}

	for i := range predictions {
		service.processPredictionScore(ctx, &predictions[i])
	}

	slog.InfoContext(ctx, "Exit: GetRoundPredictions", "user_id", userID, "year", year, "round", round, "count", len(predictions))
	return predictions, nil
}

func (service *predictionService) processPredictionScore(ctx context.Context, prediction *models.Prediction) {
	now := time.Now().UTC().UnixMilli()

	if prediction.Score != nil && prediction.RevalidateUntil != nil && now > prediction.RevalidateUntil.UnixMilli() {
		return
	}

	sessionTimeMS, err := service.getSessionStartTimeMS(ctx, prediction.Year, prediction.Round, prediction.SessionType)
	if err != nil {
		return
	}

	if prediction.RevalidateUntil == nil {
		revalidateUntil := time.UnixMilli(sessionTimeMS + (48 * 3600 * 1000)).UTC()
		prediction.RevalidateUntil = &revalidateUntil
	}

	isCompleted := now > sessionTimeMS+(2*3600*1000)
	isInRevalidationWindow := now < sessionTimeMS+(48*3600*1000)

	if !isCompleted {
		return
	}

	if prediction.Score == nil || isInRevalidationWindow {
		service.syncScoreWithResults(ctx, prediction)
	}
}

func (service *predictionService) getSessionStartTimeMS(ctx context.Context, year, round int, sessionType string) (int64, error) {
	schedule, err := service.f1Service.GetScheduleByYear(ctx, year)
	if err != nil {
		return 0, err
	}

	for _, weekend := range schedule {
		if weekend.Round == round {
			for _, session := range weekend.Sessions {
				if session.Type == sessionType {
					return session.TimeUTCMS, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("session %s not found in round %d, %d", sessionType, round, year)
}

func (service *predictionService) syncScoreWithResults(ctx context.Context, prediction *models.Prediction) {

	results, err := service.f1Service.GetSessionResults(ctx, prediction.Year, prediction.Round, prediction.SessionType)
	if err != nil || results == nil || len(results.Results) == 0 {
		return
	}

	newScore, correctness := service.scoringService.CalculateScore(prediction, results)

	if !service.hasPredictionChanged(prediction, newScore, correctness) {
		return
	}

	prediction.Score = &newScore
	if len(correctness) == len(prediction.Entries) {
		for j := range prediction.Entries {
			prediction.Entries[j].Correct = correctness[j]
		}
	} else {
		slog.WarnContext(ctx, "Scoring correctness slice length mismatch",
			"expected", len(prediction.Entries),
			"actual", len(correctness),
			"prediction_id", prediction.ID)
	}
	prediction.UpdatedAt = time.Now().UTC()

	if err := service.predictionRepository.SavePrediction(ctx, prediction); err != nil {
		slog.ErrorContext(ctx, "Failed to persist updated prediction score",
			"error", err,
			"prediction_id", prediction.ID)
	}
}

func (service *predictionService) hasPredictionChanged(prediction *models.Prediction, newScore int, correctness []bool) bool {
	if prediction.Score == nil || *prediction.Score != newScore {
		return true
	}

	if len(correctness) != len(prediction.Entries) {
		return true
	}

	for j, entry := range prediction.Entries {
		if entry.Correct != correctness[j] {
			return true
		}
	}

	return false
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
