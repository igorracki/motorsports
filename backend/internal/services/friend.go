package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/repository"
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
	var targetUser *models.User
	var err error

	if _, uuidErr := uuid.Parse(identifier); uuidErr == nil {
		targetUser, err = service.userRepo.GetUserByID(ctx, identifier)
	} else {
		targetUser, _, err = service.userRepo.GetUserByEmail(ctx, identifier)
	}

	if err != nil {
		return fmt.Errorf("searching for user: %w", err)
	}

	if targetUser == nil {
		return fmt.Errorf("%w: user %s", models.ErrNotFound, identifier)
	}

	if targetUser.ID == senderID {
		return fmt.Errorf("%w: cannot add yourself", models.ErrInvalidInput)
	}

	alreadyFriends, err := service.friendRepo.AreFriends(ctx, senderID, targetUser.ID)
	if err != nil {
		return fmt.Errorf("checking friendship status: %w", err)
	}
	if alreadyFriends {
		return fmt.Errorf("%w: already friends", models.ErrConflict)
	}

	hasPending, err := service.friendRepo.HasPendingRequest(ctx, senderID, targetUser.ID)
	if err != nil {
		return fmt.Errorf("checking pending request status: %w", err)
	}
	if hasPending {
		return fmt.Errorf("%w: request already pending", models.ErrConflict)
	}

	request := &models.FriendRequest{
		ID:         uuid.New().String(),
		SenderID:   senderID,
		ReceiverID: targetUser.ID,
		Status:     models.FriendRequestPending,
		CreatedAt:  time.Now().UTC(),
	}

	if err := service.friendRepo.CreateFriendRequest(ctx, request); err != nil {
		return fmt.Errorf("creating friend request: %w", err)
	}

	return nil
}

func (service *friendService) GetPendingRequests(ctx context.Context, userID string) ([]models.FriendRequest, error) {
	requests, err := service.friendRepo.GetPendingFriendRequestsByReceiverID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching pending requests: %w", err)
	}
	return requests, nil
}

func (service *friendService) HandleFriendRequest(ctx context.Context, userID string, requestID string, action string) error {
	request, err := service.friendRepo.GetFriendRequestByID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("fetching friend request: %w", err)
	}

	if request == nil || request.ReceiverID != userID || request.Status != models.FriendRequestPending {
		return fmt.Errorf("%w: invalid friend request", models.ErrNotFound)
	}

	if action == "accept" {
		friendship := &models.Friendship{
			UserID:    request.SenderID,
			FriendID:  request.ReceiverID,
			CreatedAt: time.Now().UTC(),
		}
		if err := service.friendRepo.AcceptFriendRequest(ctx, requestID, friendship); err != nil {
			return fmt.Errorf("accepting friend request: %w", err)
		}
	} else if action == "deny" {
		if err := service.friendRepo.UpdateFriendRequestStatus(ctx, requestID, models.FriendRequestDenied); err != nil {
			return fmt.Errorf("denying request: %w", err)
		}
	} else {
		return fmt.Errorf("%w: invalid action %s", models.ErrInvalidInput, action)
	}

	return nil
}

func (service *friendService) GetFriends(ctx context.Context, userID string) ([]string, error) {
	friends, err := service.friendRepo.GetFriendsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching friends: %w", err)
	}
	return friends, nil
}
