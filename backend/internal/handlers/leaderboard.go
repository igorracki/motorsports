package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type LeaderboardHandler struct {
	leaderboardService services.LeaderboardService
}

func NewLeaderboardHandler(leaderboardService services.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{
		leaderboardService: leaderboardService,
	}
}

func (handler *LeaderboardHandler) GetLeaderboard(context echo.Context) error {
	ctx := context.Request().Context()
	slog.InfoContext(ctx, "Entry: GetLeaderboard")

	userIDVal := context.Get("user_id")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		return context.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "user not authenticated",
		})
	}

	seasonStr := context.Param("season")
	season, err := strconv.Atoi(seasonStr)
	if err != nil {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_parameter",
			Message: "season must be a number",
		})
	}

	if season < 1950 || season > 2100 {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_parameter",
			Message: "season must be between 1950 and 2100",
		})
	}

	entries, err := handler.leaderboardService.GetLeaderboard(ctx, userID, season)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch leaderboard", "error", err, "season", season)
		return context.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "fetch_failed",
			Message: err.Error(),
		})
	}

	if entries == nil {
		entries = []models.LeaderboardEntry{}
	}

	slog.InfoContext(ctx, "Exit: GetLeaderboard", "user_id", userID, "season", season, "count", len(entries))
	return context.JSON(http.StatusOK, entries)
}
