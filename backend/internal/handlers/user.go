package handlers

import (
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

	if userID == "" {
		return models.ErrInvalidInput
	}

	profile, err := handler.userService.GetFullProfile(ctx, userID)
	if err != nil {
		return err
	}

	return context.JSON(http.StatusOK, profile)
}
