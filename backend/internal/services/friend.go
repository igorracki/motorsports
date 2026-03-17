package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/repository"
)

type FriendService interface {
	SendFriendRequest(ctx context.Context, senderID string, identifier string) error
	GetPendingRequests(ctx context.Context, userID string) ([]models.FriendRequest, error)
	HandleFriendRequest(ctx context.Context, userID string, requestID string, action string) error
	GetFriends(ctx context.Context, userID string) ([]string, error)
}

type friendService struct {
	friendRepo repository.FriendRepository
	userRepo   repository.UserRepository
}

func NewFriendService(friendRepo repository.FriendRepository, userRepo repository.UserRepository) FriendService {
	return &friendService{
		friendRepo: friendRepo,
		userRepo:   userRepo,
	}
}

func (service *friendService) SendFriendRequest(ctx context.Context, senderID string, identifier string) error {
	slog.InfoContext(ctx, "Entry: SendFriendRequest", "sender_id", senderID, "identifier", identifier)

	var targetUser *models.User
	var err error

	// Try to find by UUID first
	if _, uuidErr := uuid.Parse(identifier); uuidErr == nil {
		targetUser, err = service.userRepo.GetUserByID(ctx, identifier)
	} else {
		// Try to find by email
		targetUser, _, err = service.userRepo.GetUserByEmail(ctx, identifier)
	}

	if err != nil {
		return fmt.Errorf("searching for user: %w", err)
	}

	if targetUser == nil {
		return fmt.Errorf("user not found")
	}

	if targetUser.ID == senderID {
		return fmt.Errorf("cannot add yourself as a friend")
	}

	// Check if already friends
	alreadyFriends, err := service.friendRepo.AreFriends(ctx, senderID, targetUser.ID)
	if err != nil {
		return fmt.Errorf("checking friendship status: %w", err)
	}
	if alreadyFriends {
		return fmt.Errorf("already friends with this user")
	}

	// Create request
	request := &models.FriendRequest{
		ID:         uuid.New().String(),
		SenderID:   senderID,
		ReceiverID: targetUser.ID,
		Status:     models.FriendRequestPending,
		CreatedAt:  time.Now().UTC(),
	}

	err = service.friendRepo.CreateFriendRequest(ctx, request)
	if err != nil {
		return fmt.Errorf("creating friend request: %w", err)
	}

	slog.InfoContext(ctx, "Exit: SendFriendRequest", "sender_id", senderID, "receiver_id", targetUser.ID)
	return nil
}

func (service *friendService) GetPendingRequests(ctx context.Context, userID string) ([]models.FriendRequest, error) {
	slog.InfoContext(ctx, "Entry: GetPendingRequests", "user_id", userID)
	requests, err := service.friendRepo.GetPendingFriendRequestsByReceiverID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching pending requests: %w", err)
	}
	slog.InfoContext(ctx, "Exit: GetPendingRequests", "user_id", userID, "count", len(requests))
	return requests, nil
}

func (service *friendService) HandleFriendRequest(ctx context.Context, userID string, requestID string, action string) error {
	slog.InfoContext(ctx, "Entry: HandleFriendRequest", "user_id", userID, "request_id", requestID, "action", action)

	request, err := service.friendRepo.GetFriendRequestByID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("fetching friend request: %w", err)
	}

	if request == nil || request.ReceiverID != userID || request.Status != models.FriendRequestPending {
		return fmt.Errorf("invalid friend request")
	}

	if action == "accept" {
		err = service.friendRepo.UpdateFriendRequestStatus(ctx, requestID, models.FriendRequestAccepted)
		if err != nil {
			return fmt.Errorf("updating request status: %w", err)
		}

		friendship := &models.Friendship{
			UserID:    request.SenderID,
			FriendID:  request.ReceiverID,
			CreatedAt: time.Now().UTC(),
		}
		err = service.friendRepo.CreateFriendship(ctx, friendship)
		if err != nil {
			return fmt.Errorf("creating friendship: %w", err)
		}
	} else if action == "deny" {
		err = service.friendRepo.UpdateFriendRequestStatus(ctx, requestID, models.FriendRequestDenied)
		if err != nil {
			return fmt.Errorf("denying request: %w", err)
		}
	} else {
		return fmt.Errorf("invalid action")
	}

	slog.InfoContext(ctx, "Exit: HandleFriendRequest", "user_id", userID, "request_id", requestID)
	return nil
}

func (service *friendService) GetFriends(ctx context.Context, userID string) ([]string, error) {
	slog.InfoContext(ctx, "Entry: GetFriends", "user_id", userID)
	friends, err := service.friendRepo.GetFriendsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching friends: %w", err)
	}
	slog.InfoContext(ctx, "Exit: GetFriends", "user_id", userID, "count", len(friends))
	return friends, nil
}
