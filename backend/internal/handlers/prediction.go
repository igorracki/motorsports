package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type PredictionHandler struct {
	predictionService services.PredictionService
	scoringService    services.ScoringService
}

func NewPredictionHandler(service services.PredictionService, scoring services.ScoringService) *PredictionHandler {
	return &PredictionHandler{
		predictionService: service,
		scoringService:    scoring,
	}
}

func (handler *PredictionHandler) GetScoringRules(context echo.Context) error {
	rules := handler.scoringService.GetScoringRules()
	return context.JSON(http.StatusOK, rules)
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

	// Security: Ensure user can only submit their own predictions
	authUserID := context.Get("user_id").(string)
	if authUserID != userID {
		return context.JSON(http.StatusForbidden, models.ErrorResponse{
			Error:   "forbidden",
			Message: "cannot submit predictions for other users",
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

func (handler *PredictionHandler) GetUserPredictions(context echo.Context) error {
	ctx := context.Request().Context()
	userID := context.Param("id")
	slog.InfoContext(ctx, "Entry: GetUserPredictions", "user_id", userID)

	if userID == "" {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_parameter",
			Message: "must provide a user id in the path",
		})
	}

	// Security: Ensure user can only access their own predictions
	authUserID := context.Get("user_id").(string)
	if authUserID != userID {
		return context.JSON(http.StatusForbidden, models.ErrorResponse{
			Error:   "forbidden",
			Message: "cannot access other user predictions",
		})
	}

	predictions, err := handler.predictionService.GetUserPredictions(ctx, userID)
	if err != nil {
		slog.WarnContext(ctx, "Failed to fetch user predictions", "user_id", userID, "error", err)
		return context.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "fetch_failed",
			Message: "failed to retrieve predictions",
		})
	}

	slog.InfoContext(ctx, "Exit: GetUserPredictions", "user_id", userID, "count", len(predictions))
	return context.JSON(http.StatusOK, predictions)
}

func (handler *PredictionHandler) GetRoundPredictions(context echo.Context) error {
	ctx := context.Request().Context()
	userID := context.Param("id")
	yearParam := context.Param("year")
	roundParam := context.Param("round")

	slog.InfoContext(ctx, "Entry: GetRoundPredictions",
		"user_id", userID, "year", yearParam, "round", roundParam)

	if userID == "" || yearParam == "" || roundParam == "" {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_parameter",
			Message: "must provide user id, year, and round",
		})
	}

	// Security: Ensure user can only access their own predictions
	authUserID := context.Get("user_id").(string)
	if authUserID != userID {
		return context.JSON(http.StatusForbidden, models.ErrorResponse{
			Error:   "forbidden",
			Message: "cannot access other user predictions",
		})
	}

	var year, round int
	if _, err := fmt.Sscanf(yearParam, "%d", &year); err != nil {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_parameter",
			Message: "year must be an integer",
		})
	}
	if _, err := fmt.Sscanf(roundParam, "%d", &round); err != nil {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_parameter",
			Message: "round must be an integer",
		})
	}

	predictions, err := handler.predictionService.GetRoundPredictions(ctx, userID, year, round)
	if err != nil {
		slog.WarnContext(ctx, "Failed to fetch round predictions",
			"user_id", userID, "year", year, "round", round, "error", err)
		return context.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "fetch_failed",
			Message: "failed to retrieve predictions",
		})
	}

	slog.InfoContext(ctx, "Exit: GetRoundPredictions",
		"user_id", userID, "year", year, "round", round, "count", len(predictions))
	return context.JSON(http.StatusOK, predictions)
}
