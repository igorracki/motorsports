package middleware

import (
	"context"
	"net/http"

	"github.com/igorracki/f1/backend/internal/auth"
	f1context "github.com/igorracki/f1/backend/internal/context"
	"github.com/igorracki/f1/backend/internal/models"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("auth_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "missing authentication token",
			})
		}

		claims, err := auth.ValidateToken(cookie.Value)
		if err != nil {
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
