package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/igorracki/motorsports/backend/internal/auth"
	f1context "github.com/igorracki/motorsports/backend/internal/context"
	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("auth_token")
		if err != nil {
			slog.WarnContext(c.Request().Context(), "Authentication failure",
				"reason", "missing authentication token",
				"ip", c.RealIP(),
				"path", c.Request().URL.Path,
			)
			return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "missing authentication token",
			})
		}

		claims, err := auth.ValidateToken(cookie.Value)
		if err != nil {
			slog.WarnContext(c.Request().Context(), "Authentication failure",
				"reason", err.Error(),
				"ip", c.RealIP(),
				"path", c.Request().URL.Path,
			)
			return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "invalid or expired token",
			})
		}

		// Inject user_id into both echo context and request context
		c.Set("user_id", claims.UserID)

		ctx := context.WithValue(c.Request().Context(), f1context.UserIDKey, claims.UserID)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
