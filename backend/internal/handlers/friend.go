package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/microcosm-cc/bluemonday"
)

type FriendHandler struct {
	friendService services.FriendService
	sanitizer     *bluemonday.Policy
}

func NewFriendHandler(friendService services.FriendService) *FriendHandler {
	return &FriendHandler{
		friendService: friendService,
		sanitizer:     bluemonday.StrictPolicy(),
	}
}

func (handler *FriendHandler) SendFriendRequest(context echo.Context) error {
	ctx := context.Request().Context()
	slog.InfoContext(ctx, "Entry: SendFriendRequest")

	userIDVal := context.Get("user_id")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		return context.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "user not authenticated",
		})
	}

	var req models.SendFriendRequestRequest
	if err := context.Bind(&req); err != nil {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "invalid request body",
		})
	}

	if req.Identifier == "" {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_parameter",
			Message: "identifier is required",
		})
	}

	req.Identifier = handler.sanitizer.Sanitize(req.Identifier)

	err := handler.friendService.SendFriendRequest(ctx, userID, req.Identifier)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to send friend request", "error", err)

		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrUserNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, services.ErrCannotAddSelf) || errors.Is(err, services.ErrAlreadyFriends) || errors.Is(err, services.ErrRequestAlreadyPending) {
			status = http.StatusBadRequest
		}

		return context.JSON(status, models.ErrorResponse{
			Error:   "friend_request_failed",
			Message: err.Error(),
		})
	}

	slog.InfoContext(ctx, "Exit: SendFriendRequest", "user_id", userID)
	return context.NoContent(http.StatusCreated)
}

func (handler *FriendHandler) GetPendingRequests(context echo.Context) error {
	ctx := context.Request().Context()
	slog.InfoContext(ctx, "Entry: GetPendingRequests")

	userIDVal := context.Get("user_id")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		return context.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "user not authenticated",
		})
	}

	requests, err := handler.friendService.GetPendingRequests(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch pending requests", "error", err)
		return context.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "fetch_failed",
			Message: err.Error(),
		})
	}

	if requests == nil {
		requests = []models.FriendRequest{}
	}

	slog.InfoContext(ctx, "Exit: GetPendingRequests", "user_id", userID, "count", len(requests))
	return context.JSON(http.StatusOK, requests)
}

func (handler *FriendHandler) HandleFriendRequest(context echo.Context) error {
	ctx := context.Request().Context()
	slog.InfoContext(ctx, "Entry: HandleFriendRequest")

	userIDVal := context.Get("user_id")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		return context.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "user not authenticated",
		})
	}
	requestID := context.Param("id")
	if requestID == "" {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_parameter",
			Message: "request id is required",
		})
	}

	var req models.HandleFriendRequestRequest
	if err := context.Bind(&req); err != nil {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "invalid request body",
		})
	}

	req.Action = handler.sanitizer.Sanitize(req.Action)

	err := handler.friendService.HandleFriendRequest(ctx, userID, requestID, req.Action)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to handle friend request", "error", err, "request_id", requestID)

		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrInvalidFriendRequest) || errors.Is(err, services.ErrInvalidAction) {
			status = http.StatusBadRequest
		}

		return context.JSON(status, models.ErrorResponse{
			Error:   "handle_failed",
			Message: err.Error(),
		})
	}

	slog.InfoContext(ctx, "Exit: HandleFriendRequest", "user_id", userID, "request_id", requestID)
	return context.NoContent(http.StatusOK)
}
