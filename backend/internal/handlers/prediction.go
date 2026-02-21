package handlers

import (
	"log/slog"
	"net/http"

	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type PredictionHandler struct {
	predictionService services.PredictionService
}

func NewPredictionHandler(service services.PredictionService) *PredictionHandler {
	return &PredictionHandler{
		predictionService: service,
	}
}

func (handler *PredictionHandler) SubmitPrediction(context echo.Context) error {
	ctx := context.Request().Context()
	userID := context.Param("id")
	slog.InfoContext(ctx, "Entry: SubmitPrediction", "user_id", userID)

	if userID == "" {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_parameter",
			Message: "must provide a user id in the path",
		})
	}

	var prediction models.Prediction
	if err := context.Bind(&prediction); err != nil {
		slog.WarnContext(ctx, "Failed to bind prediction request", "error", err)
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "failed to parse request body",
		})
	}

	// Override or set the UserID from the URL parameter for security/consistency
	prediction.UserID = userID

	if err := handler.predictionService.SubmitPrediction(ctx, &prediction); err != nil {
		slog.WarnContext(ctx, "Prediction submission failed",
			"user_id", prediction.UserID,
			"error", err)
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "submission_failed",
			Message: err.Error(),
		})
	}

	slog.InfoContext(ctx, "Exit: SubmitPrediction",
		"user_id", prediction.UserID,
		"prediction_id", prediction.ID)
	return context.JSON(http.StatusCreated, prediction)
}
