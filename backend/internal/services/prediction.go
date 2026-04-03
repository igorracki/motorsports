package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/repository"
)

type PredictionService interface {
	SubmitPrediction(ctx context.Context, prediction *models.Prediction) error
	GetUserPredictions(ctx context.Context, userID string) ([]models.Prediction, error)
	GetRoundPredictions(ctx context.Context, userID string, year, round int) ([]models.Prediction, error)
	GetPolicyConfig() models.PredictionPolicyConfig
}

type predictionService struct {
	predictionRepository repository.PredictionRepository
	f1Service            F1Service
	scoringService       ScoringService
	policy               PredictionPolicy
	configService        ConfigService
}

func NewPredictionService(repo repository.PredictionRepository, f1 F1Service, scoring ScoringService, policy PredictionPolicy, configService ConfigService) PredictionService {
	return &predictionService{
		predictionRepository: repo,
		f1Service:            f1,
		scoringService:       scoring,
		policy:               policy,
		configService:        configService,
	}
}

func (service *predictionService) SubmitPrediction(ctx context.Context, prediction *models.Prediction) error {
	if err := service.validatePrediction(prediction); err != nil {
		return err
	}

	sessionTimeMS, err := service.getSessionStartTimeMS(ctx, prediction.Year, prediction.Round, prediction.SessionType)
	if err != nil {
		return fmt.Errorf("getting session start time: %w", err)
	}

	if service.policy.IsLocked(sessionTimeMS) {
		return fmt.Errorf("%w: prediction period has closed", models.ErrForbidden)
	}

	now := time.Now().UTC()
	revalidateUntil := service.policy.GetRevalidationDeadline(sessionTimeMS)

	prediction.ID = uuid.New().String()
	prediction.RevalidateUntil = &revalidateUntil
	prediction.CreatedAt = now
	prediction.UpdatedAt = now

	if err := service.predictionRepository.SavePrediction(ctx, prediction); err != nil {
		return fmt.Errorf("saving prediction: %w", err)
	}

	return nil
}

func (service *predictionService) GetUserPredictions(ctx context.Context, userID string) ([]models.Prediction, error) {
	predictions, err := service.predictionRepository.GetUserPredictions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching predictions: %w", err)
	}

	for i := range predictions {
		service.processPredictionScore(ctx, &predictions[i])
	}

	return predictions, nil
}

func (service *predictionService) GetRoundPredictions(ctx context.Context, userID string, year, round int) ([]models.Prediction, error) {
	predictions, err := service.predictionRepository.GetRoundPredictions(ctx, userID, year, round)
	if err != nil {
		return nil, fmt.Errorf("fetching round predictions: %w", err)
	}

	for i := range predictions {
		service.processPredictionScore(ctx, &predictions[i])
	}

	return predictions, nil
}

func (service *predictionService) GetPolicyConfig() models.PredictionPolicyConfig {
	return service.policy.GetConfig()
}

func (service *predictionService) processPredictionScore(ctx context.Context, prediction *models.Prediction) {
	sessionTimeMS, err := service.getSessionStartTimeMS(ctx, prediction.Year, prediction.Round, prediction.SessionType)
	if err != nil {
		return
	}

	if prediction.Score != nil && service.policy.IsScoringFinal(sessionTimeMS) {
		return
	}

	if prediction.RevalidateUntil == nil {
		revalidateUntil := service.policy.GetRevalidationDeadline(sessionTimeMS)
		prediction.RevalidateUntil = &revalidateUntil
	}

	// We only sync if the session is started and enough time has passed
	if !service.policy.IsLocked(sessionTimeMS) {
		return
	}

	if prediction.Score == nil || !service.policy.IsScoringFinal(sessionTimeMS) {
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
				if session.Type == sessionType || session.SessionCode == sessionType {
					return session.TimeUTCMS, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("%w: session %s not found", models.ErrNotFound, sessionType)
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
	config := service.configService.GetAppConfig().Validation
	entryCount := len(prediction.Entries)

	if entryCount < config.MinEntries || entryCount > config.MaxEntries {
		return fmt.Errorf("%w: prediction must have between %d and %d entries, got %d",
			models.ErrInvalidInput, config.MinEntries, config.MaxEntries, entryCount)
	}

	positions := make(map[int]bool)
	drivers := make(map[string]bool)

	for _, entry := range prediction.Entries {
		if entry.Position < 1 || entry.Position > config.MaxEntries {
			return fmt.Errorf("%w: invalid position %d", models.ErrInvalidInput, entry.Position)
		}
		if positions[entry.Position] {
			return fmt.Errorf("%w: duplicate position %d", models.ErrInvalidInput, entry.Position)
		}
		positions[entry.Position] = true

		if entry.DriverID == "" {
			return fmt.Errorf("%w: driver_id cannot be empty", models.ErrInvalidInput)
		}
		if drivers[entry.DriverID] {
			return fmt.Errorf("%w: duplicate driver %s", models.ErrInvalidInput, entry.DriverID)
		}
		drivers[entry.DriverID] = true
	}

	return nil
}
