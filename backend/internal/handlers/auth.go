package handlers

import (
	"log/slog"
	"net/http"
	"time"

	f1context "github.com/igorracki/f1/backend/internal/context"
	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService services.AuthService
	userService services.UserService
}

func NewAuthHandler(authService services.AuthService, userService services.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

func (handler *AuthHandler) Register(context echo.Context) error {
	ctx := context.Request().Context()
	slog.InfoContext(ctx, "Entry: Register")

	var request models.RegisterUserRequest
	if err := context.Bind(&request); err != nil {
		slog.WarnContext(ctx, "Failed to bind registration request", "error", err)
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "failed to parse request body",
		})
	}

	user, profile, err := handler.authService.Register(ctx, request)
	if err != nil {
		slog.WarnContext(ctx, "User registration failed", "email", request.Email, "error", err)
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "registration_failed",
			Message: err.Error(),
		})
	}

	slog.InfoContext(ctx, "Exit: Register", "user_id", user.ID)
	return context.JSON(http.StatusCreated, models.AuthResponse{
		User:    *user,
		Profile: *profile,
	})
}

func (handler *AuthHandler) Login(context echo.Context) error {
	ctx := context.Request().Context()
	slog.InfoContext(ctx, "Entry: Login")

	var request models.LoginRequest
	if err := context.Bind(&request); err != nil {
		slog.WarnContext(ctx, "Failed to bind login request", "error", err)
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "failed to parse request body",
		})
	}

	user, profile, token, expiresAt, err := handler.authService.Login(ctx, request)
	if err != nil {
		slog.WarnContext(ctx, "User login failed", "email", request.Email, "error", err)
		return context.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "login_failed",
			Message: err.Error(),
		})
	}

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	context.SetCookie(cookie)

	slog.InfoContext(ctx, "Exit: Login", "user_id", user.ID)
	return context.JSON(http.StatusOK, models.AuthResponse{
		User:    *user,
		Profile: *profile,
	})
}

func (handler *AuthHandler) Logout(context echo.Context) error {
	ctx := context.Request().Context()
	slog.InfoContext(ctx, "Entry: Logout")

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	context.SetCookie(cookie)

	slog.InfoContext(ctx, "Exit: Logout")
	return context.NoContent(http.StatusOK)
}

func (handler *AuthHandler) Me(context echo.Context) error {
	ctx := context.Request().Context()
	slog.InfoContext(ctx, "Entry: Me")

	userID, ok := ctx.Value(f1context.UserIDKey).(string)
	if !ok || userID == "" {
		return context.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "user not authenticated",
		})
	}

	profileResponse, err := handler.userService.GetFullProfile(ctx, userID)
	if err != nil {
		slog.WarnContext(ctx, "Failed to fetch user profile in Me", "user_id", userID, "error", err)
		return context.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "failed to fetch user profile",
		})
	}

	slog.InfoContext(ctx, "Exit: Me", "user_id", userID)
	return context.JSON(http.StatusOK, models.AuthResponse{
		User:    profileResponse.User,
		Profile: profileResponse.Profile,
	})
}
