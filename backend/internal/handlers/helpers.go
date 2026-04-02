package handlers

import (
	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/labstack/echo/v4"
)

func GetAuthenticatedUserID(context echo.Context) (string, error) {
	authUserIDVal := context.Get("user_id")
	authUserID, ok := authUserIDVal.(string)
	if !ok || authUserID == "" {
		return "", models.ErrUnauthorized
	}
	return authUserID, nil
}
