package handlers

import (
	"net/http"
	"time"

	_context "github.com/igorracki/motorsports/backend/internal/context"
	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService  services.AuthService
	userService  services.UserService
	cookieSecure bool
}

type AuthHandlerOption func(*AuthHandler)

func NewAuthHandler(authService services.AuthService, userService services.UserService, options ...AuthHandlerOption) *AuthHandler {
	handler := &AuthHandler{
		authService: authService,
		userService: userService,
	}

	for _, option := range options {
		option(handler)
	}

	return handler
}

func WithCookieSecure(secure bool) AuthHandlerOption {
	return func(handler *AuthHandler) {
		handler.cookieSecure = secure
	}
}

func (handler *AuthHandler) Register(context echo.Context) error {
	ctx := context.Request().Context()

	var request models.RegisterUserRequest
	if err := context.Bind(&request); err != nil {
		return models.ErrInvalidInput
	}

	if err := context.Validate(&request); err != nil {
		return models.ErrInvalidInput
	}

	user, profile, token, expiresAt, err := handler.authService.Register(ctx, request)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: true,
		Secure:   handler.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
	context.SetCookie(cookie)

	return context.JSON(http.StatusCreated, models.AuthResponse{
		User:    *user,
		Profile: *profile,
	})
}

func (handler *AuthHandler) Login(context echo.Context) error {
	ctx := context.Request().Context()

	var request models.LoginRequest
	if err := context.Bind(&request); err != nil {
		return models.ErrInvalidInput
	}

	user, profile, token, expiresAt, err := handler.authService.Login(ctx, request)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: true,
		Secure:   handler.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
	context.SetCookie(cookie)

	return context.JSON(http.StatusOK, models.AuthResponse{
		User:    *user,
		Profile: *profile,
	})
}

func (handler *AuthHandler) Logout(context echo.Context) error {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   handler.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
	context.SetCookie(cookie)
	return context.NoContent(http.StatusNoContent)
}

func (handler *AuthHandler) Me(context echo.Context) error {
	ctx := context.Request().Context()

	userID, ok := ctx.Value(_context.UserIDKey).(string)
	if !ok || userID == "" {
		return models.ErrUnauthorized
	}

	profileResponse, err := handler.userService.GetFullProfile(ctx, userID)
	if err != nil {
		return err
	}

	return context.JSON(http.StatusOK, models.AuthResponse{
		User:    profileResponse.User,
		Profile: profileResponse.Profile,
	})
}
