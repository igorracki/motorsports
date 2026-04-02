package handlers

import (
	"net/http"

	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type FriendHandler struct {
	friendService services.FriendService
}

func NewFriendHandler(friendService services.FriendService) *FriendHandler {
	return &FriendHandler{
		friendService: friendService,
	}
}

func (handler *FriendHandler) SendFriendRequest(context echo.Context) error {
	ctx := context.Request().Context()

	userID, err := GetAuthenticatedUserID(context)
	if err != nil {
		return err
	}

	var request models.SendFriendRequestRequest
	if err := context.Bind(&request); err != nil {
		return models.ErrInvalidInput
	}
	if err := context.Validate(&request); err != nil {
		return models.ErrInvalidInput
	}

	if err := handler.friendService.SendFriendRequest(ctx, userID, request.Identifier); err != nil {
		return err
	}

	return context.NoContent(http.StatusCreated)
}

func (handler *FriendHandler) GetPendingRequests(context echo.Context) error {
	ctx := context.Request().Context()

	userID, err := GetAuthenticatedUserID(context)
	if err != nil {
		return err
	}

	requests, err := handler.friendService.GetPendingRequests(ctx, userID)
	if err != nil {
		return err
	}

	if requests == nil {
		requests = []models.FriendRequest{}
	}

	return context.JSON(http.StatusOK, requests)
}

func (handler *FriendHandler) HandleFriendRequest(context echo.Context) error {
	ctx := context.Request().Context()

	userID, err := GetAuthenticatedUserID(context)
	if err != nil {
		return err
	}

	requestID := context.Param("id")
	if requestID == "" {
		return models.ErrInvalidInput
	}

	var request models.HandleFriendRequestRequest
	if err := context.Bind(&request); err != nil {
		return models.ErrInvalidInput
	}
	if err := context.Validate(&request); err != nil {
		return models.ErrInvalidInput
	}

	if err := handler.friendService.HandleFriendRequest(ctx, userID, requestID, request.Action); err != nil {
		return err
	}

	return context.NoContent(http.StatusOK)
}
