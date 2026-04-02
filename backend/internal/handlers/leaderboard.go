package handlers

import (
	"net/http"

	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/services"
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

type LeaderboardParam struct {
	Season int `param:"season" validate:"required,min=1950,max=2100"`
}

func (handler *LeaderboardHandler) GetLeaderboard(context echo.Context) error {
	ctx := context.Request().Context()

	userID, err := GetAuthenticatedUserID(context)
	if err != nil {
		return err
	}

	var params LeaderboardParam
	if err := context.Bind(&params); err != nil {
		return models.ErrInvalidInput
	}
	if err := context.Validate(&params); err != nil {
		return models.ErrInvalidInput
	}

	entries, err := handler.leaderboardService.GetLeaderboard(ctx, userID, params.Season)
	if err != nil {
		return err
	}

	if entries == nil {
		entries = []models.LeaderboardEntry{}
	}

	return context.JSON(http.StatusOK, entries)
}
