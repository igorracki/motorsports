package middleware

import (
	"context"
	"net/http"

	"github.com/igorracki/motorsports/backend/internal/auth"
	_context "github.com/igorracki/motorsports/backend/internal/context"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(tokenManager auth.TokenManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(contextObj echo.Context) error {
			cookie, err := contextObj.Cookie("auth_token")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authentication token")
			}

			claims, err := tokenManager.ValidateToken(cookie.Value)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			contextObj.Set("user_id", claims.UserID)

			ctx := context.WithValue(contextObj.Request().Context(), _context.UserIDKey, claims.UserID)
			contextObj.SetRequest(contextObj.Request().WithContext(ctx))

			return next(contextObj)
		}
	}
}

// RequireResourceOwnerMiddleware ensures that the authenticated user matches the "id" path parameter.
func RequireResourceOwnerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authUserIDVal := c.Get("user_id")
			authUserID, ok := authUserIDVal.(string)
			if !ok || authUserID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
			}

			resourceID := c.Param("id")
			if resourceID != "" && resourceID != authUserID {
				return echo.NewHTTPError(http.StatusForbidden, "permission denied")
			}

			return next(c)
		}
	}
}
