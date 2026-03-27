package handlers

import (
	"log/slog"
	"net/http"

	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (handler *UserHandler) GetUserProfile(context echo.Context) error {
	ctx := context.Request().Context()
	userID := context.Param("id")
	slog.InfoContext(ctx, "Entry: GetUserProfile", "user_id", userID)

	if userID == "" {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_parameter",
			Message: "must provide a user id",
		})
	}

	// Security: Ensure user can only access their own profile
	authUserIDVal := context.Get("user_id")
	authUserID, ok := authUserIDVal.(string)
	if !ok || authUserID == "" {
		return context.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "authentication required",
		})
	}
	if authUserID != userID {
		return context.JSON(http.StatusForbidden, models.ErrorResponse{
			Error:   "forbidden",
			Message: "cannot access other user profiles",
		})
	}

	profile, err := handler.userService.GetFullProfile(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch user profile", "user_id", userID, "error", err)
		return context.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "profile_not_found",
			Message: "failed to retrieve user profile",
		})
	}

	slog.InfoContext(ctx, "Exit: GetUserProfile", "user_id", userID)
	return context.JSON(http.StatusOK, profile)
}

func (handler *UserHandler) GetSeasonScores(context echo.Context) error {
	ctx := context.Request().Context()
	userID := context.Param("id")
	slog.InfoContext(ctx, "Entry: GetSeasonScores", "user_id", userID)

	if userID == "" {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_parameter",
			Message: "must provide a user id",
		})
	}

	// Security: Ensure user can only access their own stats
	authUserIDVal := context.Get("user_id")
	authUserID, ok := authUserIDVal.(string)
	if !ok || authUserID == "" {
		return context.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "authentication required",
		})
	}
	if authUserID != userID {
		return context.JSON(http.StatusForbidden, models.ErrorResponse{
			Error:   "forbidden",
			Message: "cannot access other user stats",
		})
	}

	scores, err := handler.userService.GetSeasonScores(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch season scores", "user_id", userID, "error", err)
		return context.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "fetch_failed",
			Message: "failed to retrieve season scores",
		})
	}

	if scores == nil {
		scores = []models.UserScore{}
	}

	slog.InfoContext(ctx, "Exit: GetSeasonScores", "user_id", userID, "count", len(scores))
	return context.JSON(http.StatusOK, scores)
}
