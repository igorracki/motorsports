package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/igorracki/motorsports/backend/internal/database"
	"github.com/igorracki/motorsports/backend/internal/models"
)

const (
	insertFriendRequestSQL      = "INSERT INTO friend_requests (id, sender_id, receiver_id, status, created_at) VALUES (?, ?, ?, ?, ?)"
	getFriendRequestByIDSQL     = "SELECT id, sender_id, receiver_id, status, created_at FROM friend_requests WHERE id = ?"
	getPendingFriendRequestsSQL = `
		SELECT fr.id, fr.sender_id, fr.receiver_id, fr.status, fr.created_at, u.email, p.display_name
		FROM friend_requests fr
		JOIN users u ON fr.sender_id = u.id
		JOIN profiles p ON fr.sender_id = p.user_id
		WHERE fr.receiver_id = ? AND fr.status = ?`
	updateFriendRequestStatusSQL = "UPDATE friend_requests SET status = ? WHERE id = ?"
	insertFriendshipSQL          = "INSERT INTO friendships (user_id, friend_id, created_at) VALUES (?, ?, ?)"
	getFriendsByUserIDSQL        = "SELECT friend_id FROM friendships WHERE user_id = ?"
	checkFriendshipExistsSQL     = "SELECT EXISTS(SELECT 1 FROM friendships WHERE user_id = ? AND friend_id = ?)"
	checkPendingRequestExistsSQL = "SELECT EXISTS(SELECT 1 FROM friend_requests WHERE ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) AND status = ?)"
)

type FriendRepository interface {
	CreateFriendRequest(ctx context.Context, request *models.FriendRequest) error
	GetFriendRequestByID(ctx context.Context, id string) (*models.FriendRequest, error)
	GetPendingFriendRequestsByReceiverID(ctx context.Context, receiverID string) ([]models.FriendRequest, error)
	UpdateFriendRequestStatus(ctx context.Context, id string, status models.FriendRequestStatus) error
	CreateFriendship(ctx context.Context, friendship *models.Friendship) error
	GetFriendsByUserID(ctx context.Context, userID string) ([]string, error)
	AreFriends(ctx context.Context, user1ID, user2ID string) (bool, error)
	HasPendingRequest(ctx context.Context, user1ID, user2ID string) (bool, error)
	AcceptFriendRequest(ctx context.Context, requestID string, friendship *models.Friendship) error
}

type friendRepository struct {
	manager *database.Manager
}

func NewFriendRepository(manager *database.Manager) FriendRepository {
	return &friendRepository{manager: manager}
}

func (repo *friendRepository) CreateFriendRequest(ctx context.Context, request *models.FriendRequest) error {
	_, err := repo.manager.DB().ExecContext(ctx, insertFriendRequestSQL, request.ID, request.SenderID, request.ReceiverID, request.Status, request.CreatedAt)
	if err != nil {
		return fmt.Errorf("inserting friend request: %w", err)
	}
	return nil
}

func (repo *friendRepository) GetFriendRequestByID(ctx context.Context, id string) (*models.FriendRequest, error) {
	request := &models.FriendRequest{}
	err := repo.manager.DB().QueryRowContext(ctx, getFriendRequestByIDSQL, id).Scan(&request.ID, &request.SenderID, &request.ReceiverID, &request.Status, &request.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("querying friend request %s: %w", id, err)
	}

	return request, nil
}

func (repo *friendRepository) GetPendingFriendRequestsByReceiverID(ctx context.Context, receiverID string) ([]models.FriendRequest, error) {
	rows, err := repo.manager.DB().QueryContext(ctx, getPendingFriendRequestsSQL, receiverID, models.FriendRequestPending)
	if err != nil {
		return nil, fmt.Errorf("querying friend requests: %w", err)
	}
	defer rows.Close()

	requests := make([]models.FriendRequest, 0)
	for rows.Next() {
		var req models.FriendRequest
		if err := rows.Scan(&req.ID, &req.SenderID, &req.ReceiverID, &req.Status, &req.CreatedAt, &req.SenderEmail, &req.SenderName); err != nil {
			return nil, fmt.Errorf("scanning friend request: %w", err)
		}
		requests = append(requests, req)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating friend requests: %w", err)
	}

	return requests, nil
}

func (repo *friendRepository) UpdateFriendRequestStatus(ctx context.Context, id string, status models.FriendRequestStatus) error {
	_, err := repo.manager.DB().ExecContext(ctx, updateFriendRequestStatusSQL, status, id)
	if err != nil {
		return fmt.Errorf("updating friend request status: %w", err)
	}
	return nil
}

func (repo *friendRepository) CreateFriendship(ctx context.Context, friendship *models.Friendship) error {
	return repo.manager.Transaction(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, insertFriendshipSQL, friendship.UserID, friendship.FriendID, friendship.CreatedAt); err != nil {
			return fmt.Errorf("inserting friendship direction 1: %w", err)
		}
		if _, err := tx.ExecContext(ctx, insertFriendshipSQL, friendship.FriendID, friendship.UserID, friendship.CreatedAt); err != nil {
			return fmt.Errorf("inserting friendship direction 2: %w", err)
		}
		return nil
	})
}

func (repo *friendRepository) GetFriendsByUserID(ctx context.Context, userID string) ([]string, error) {
	rows, err := repo.manager.DB().QueryContext(ctx, getFriendsByUserIDSQL, userID)
	if err != nil {
		return nil, fmt.Errorf("querying friendships: %w", err)
	}
	defer rows.Close()

	friendIDs := make([]string, 0)
	for rows.Next() {
		var friendID string
		if err := rows.Scan(&friendID); err != nil {
			return nil, fmt.Errorf("scanning friend id: %w", err)
		}
		friendIDs = append(friendIDs, friendID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating friendships: %w", err)
	}

	return friendIDs, nil
}

func (repo *friendRepository) AreFriends(ctx context.Context, user1ID, user2ID string) (bool, error) {
	var exists bool
	err := repo.manager.DB().QueryRowContext(ctx, checkFriendshipExistsSQL, user1ID, user2ID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking friendship: %w", err)
	}
	return exists, nil
}

func (repo *friendRepository) HasPendingRequest(ctx context.Context, user1ID, user2ID string) (bool, error) {
	var exists bool
	err := repo.manager.DB().QueryRowContext(ctx, checkPendingRequestExistsSQL, user1ID, user2ID, user2ID, user1ID, models.FriendRequestPending).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking pending request: %w", err)
	}
	return exists, nil
}

func (repo *friendRepository) AcceptFriendRequest(ctx context.Context, requestID string, friendship *models.Friendship) error {
	return repo.manager.Transaction(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, updateFriendRequestStatusSQL, models.FriendRequestAccepted, requestID); err != nil {
			return fmt.Errorf("updating friend request status: %w", err)
		}
		if _, err := tx.ExecContext(ctx, insertFriendshipSQL, friendship.UserID, friendship.FriendID, friendship.CreatedAt); err != nil {
			return fmt.Errorf("inserting friendship direction 1: %w", err)
		}
		if _, err := tx.ExecContext(ctx, insertFriendshipSQL, friendship.FriendID, friendship.UserID, friendship.CreatedAt); err != nil {
			return fmt.Errorf("inserting friendship direction 2: %w", err)
		}
		return nil
	})
}
