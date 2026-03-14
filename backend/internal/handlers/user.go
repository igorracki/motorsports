package handlers

import (
	"log/slog"
	"net/http"

	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/services"
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

	profile, err := handler.userService.GetFullProfile(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch user profile", "user_id", userID, "error", err)
		return context.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "profile_not_found",
			Message: err.Error(),
		})
	}

	slog.InfoContext(ctx, "Exit: GetUserProfile", "user_id", userID)
	return context.JSON(http.StatusOK, profile)
}
