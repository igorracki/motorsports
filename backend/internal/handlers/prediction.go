package handlers

import (
	"fmt"
	"net/http"

	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type PredictionHandler struct {
	predictionService services.PredictionService
	scoringService    services.ScoringService
	configService     services.ConfigService
}

type GetRoundPredictionsParams struct {
	ID    string `param:"id" validate:"required"`
	Year  int    `param:"year" validate:"required,min=1950,max=2100"`
	Round int    `param:"round" validate:"required,min=1,max=50"`
}

func NewPredictionHandler(service services.PredictionService, scoring services.ScoringService, config services.ConfigService) *PredictionHandler {
	return &PredictionHandler{
		predictionService: service,
		scoringService:    scoring,
		configService:     config,
	}
}

func (handler *PredictionHandler) GetScoringRules(context echo.Context) error {
	rules := handler.scoringService.GetScoringRules()
	return context.JSON(http.StatusOK, rules)
}

func (handler *PredictionHandler) GetPredictionPolicy(context echo.Context) error {
	policy := handler.predictionService.GetPolicyConfig()
	return context.JSON(http.StatusOK, policy)
}

func (handler *PredictionHandler) SubmitPrediction(context echo.Context) error {
	ctx := context.Request().Context()
	userID := context.Param("id")

	var prediction models.Prediction
	if err := context.Bind(&prediction); err != nil {
		return models.ErrInvalidInput
	}

	if err := context.Validate(&prediction); err != nil {
		return models.ErrInvalidInput
	}

	prediction.UserID = userID

	if err := handler.validatePrediction(&prediction); err != nil {
		return err
	}

	if err := handler.predictionService.SubmitPrediction(ctx, &prediction); err != nil {
		return err
	}

	return context.JSON(http.StatusCreated, prediction)
}

func (handler *PredictionHandler) GetRoundPredictions(context echo.Context) error {
	ctx := context.Request().Context()

	var params GetRoundPredictionsParams
	if err := context.Bind(&params); err != nil {
		return models.ErrInvalidInput
	}
	if err := context.Validate(&params); err != nil {
		return models.ErrInvalidInput
	}

	predictions, err := handler.predictionService.GetRoundPredictions(ctx, params.ID, params.Year, params.Round)
	if err != nil {
		return err
	}

	if predictions == nil {
		predictions = []models.Prediction{}
	}

	return context.JSON(http.StatusOK, predictions)
}

func (handler *PredictionHandler) validatePrediction(prediction *models.Prediction) error {
	config := handler.configService.GetAppConfig().Validation
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
