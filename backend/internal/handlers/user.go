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

func (handler *UserHandler) RegisterUser(context echo.Context) error {
	ctx := context.Request().Context()
	slog.InfoContext(ctx, "Entry: RegisterUser")

	var request models.RegisterUserRequest
	if err := context.Bind(&request); err != nil {
		slog.WarnContext(ctx, "Failed to bind registration request", "error", err)
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "failed to parse request body",
		})
	}

	user, err := handler.userService.RegisterUser(ctx, request)
	if err != nil {
		slog.WarnContext(ctx, "User registration failed", "username", request.Username, "error", err)
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "registration_failed",
			Message: err.Error(),
		})
	}

	slog.InfoContext(ctx, "Exit: RegisterUser", "user_id", user.ID, "username", user.Username)
	return context.JSON(http.StatusCreated, user)
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
