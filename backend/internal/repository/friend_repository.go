package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/igorracki/motorsports/backend/internal/models"
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
	database *sql.DB
}

func NewFriendRepository(db *sql.DB) FriendRepository {
	return &friendRepository{database: db}
}

func (repo *friendRepository) CreateFriendRequest(ctx context.Context, request *models.FriendRequest) error {
	slog.InfoContext(ctx, "Entry: CreateFriendRequest", "sender_id", request.SenderID, "receiver_id", request.ReceiverID)

	_, err := repo.database.ExecContext(ctx,
		"INSERT INTO friend_requests (id, sender_id, receiver_id, status, created_at) VALUES (?, ?, ?, ?, ?)",
		request.ID, request.SenderID, request.ReceiverID, request.Status, request.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting friend request: %w", err)
	}

	slog.InfoContext(ctx, "Exit: CreateFriendRequest", "sender_id", request.SenderID, "receiver_id", request.ReceiverID)
	return nil
}

func (repo *friendRepository) GetFriendRequestByID(ctx context.Context, id string) (*models.FriendRequest, error) {
	slog.InfoContext(ctx, "Entry: GetFriendRequestByID", "request_id", id)
	request := &models.FriendRequest{}
	err := repo.database.QueryRowContext(ctx,
		"SELECT id, sender_id, receiver_id, status, created_at FROM friend_requests WHERE id = ?",
		id,
	).Scan(&request.ID, &request.SenderID, &request.ReceiverID, &request.Status, &request.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.InfoContext(ctx, "Exit: GetFriendRequestByID", "request_id", id, "found", false)
			return nil, nil
		}
		return nil, fmt.Errorf("querying friend request %s: %w", id, err)
	}

	slog.InfoContext(ctx, "Exit: GetFriendRequestByID", "request_id", id, "found", true)
	return request, nil
}

func (repo *friendRepository) GetPendingFriendRequestsByReceiverID(ctx context.Context, receiverID string) ([]models.FriendRequest, error) {
	slog.InfoContext(ctx, "Entry: GetPendingFriendRequestsByReceiverID", "receiver_id", receiverID)

	query := `
		SELECT fr.id, fr.sender_id, fr.receiver_id, fr.status, fr.created_at, u.email, p.display_name
		FROM friend_requests fr
		JOIN users u ON fr.sender_id = u.id
		JOIN profiles p ON fr.sender_id = p.user_id
		WHERE fr.receiver_id = ? AND fr.status = ?
	`
	rows, err := repo.database.QueryContext(ctx, query, receiverID, models.FriendRequestPending)
	if err != nil {
		return nil, fmt.Errorf("querying friend requests: %w", err)
	}
	defer rows.Close()

	requests := []models.FriendRequest{}
	for rows.Next() {
		var req models.FriendRequest
		err := rows.Scan(&req.ID, &req.SenderID, &req.ReceiverID, &req.Status, &req.CreatedAt, &req.SenderEmail, &req.SenderName)
		if err != nil {
			return nil, fmt.Errorf("scanning friend request: %w", err)
		}
		requests = append(requests, req)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating friend requests: %w", err)
	}

	slog.InfoContext(ctx, "Exit: GetPendingFriendRequestsByReceiverID", "receiver_id", receiverID, "count", len(requests))
	return requests, nil
}

func (repo *friendRepository) UpdateFriendRequestStatus(ctx context.Context, id string, status models.FriendRequestStatus) error {
	slog.InfoContext(ctx, "Entry: UpdateFriendRequestStatus", "request_id", id, "status", string(status))

	_, err := repo.database.ExecContext(ctx,
		"UPDATE friend_requests SET status = ? WHERE id = ?",
		status, id,
	)
	if err != nil {
		return fmt.Errorf("updating friend request status: %w", err)
	}

	slog.InfoContext(ctx, "Exit: UpdateFriendRequestStatus", "request_id", id)
	return nil
}

func (repo *friendRepository) CreateFriendship(ctx context.Context, friendship *models.Friendship) error {
	slog.InfoContext(ctx, "Entry: CreateFriendship", "user_id", friendship.UserID, "friend_id", friendship.FriendID)

	transaction, err := repo.database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer transaction.Rollback()

	// Insert both directions
	query := "INSERT INTO friendships (user_id, friend_id, created_at) VALUES (?, ?, ?)"

	_, err = transaction.ExecContext(ctx, query, friendship.UserID, friendship.FriendID, friendship.CreatedAt)
	if err != nil {
		return fmt.Errorf("inserting friendship direction 1: %w", err)
	}

	_, err = transaction.ExecContext(ctx, query, friendship.FriendID, friendship.UserID, friendship.CreatedAt)
	if err != nil {
		return fmt.Errorf("inserting friendship direction 2: %w", err)
	}

	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	slog.InfoContext(ctx, "Exit: CreateFriendship", "user_id", friendship.UserID, "friend_id", friendship.FriendID)
	return nil
}

func (repo *friendRepository) GetFriendsByUserID(ctx context.Context, userID string) ([]string, error) {
	slog.InfoContext(ctx, "Entry: GetFriendsByUserID", "user_id", userID)
	rows, err := repo.database.QueryContext(ctx,
		"SELECT friend_id FROM friendships WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying friendships: %w", err)
	}
	defer rows.Close()

	friendIDs := []string{}
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

	slog.InfoContext(ctx, "Exit: GetFriendsByUserID", "user_id", userID, "count", len(friendIDs))
	return friendIDs, nil
}

func (repo *friendRepository) AreFriends(ctx context.Context, user1ID, user2ID string) (bool, error) {
	slog.InfoContext(ctx, "Entry: AreFriends", "user1_id", user1ID, "user2_id", user2ID)
	var exists bool
	err := repo.database.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM friendships WHERE user_id = ? AND friend_id = ?)",
		user1ID, user2ID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("checking friendship: %w", err)
	}

	slog.InfoContext(ctx, "Exit: AreFriends", "user1_id", user1ID, "user2_id", user2ID, "are_friends", exists)
	return exists, nil
}

func (repo *friendRepository) HasPendingRequest(ctx context.Context, user1ID, user2ID string) (bool, error) {
	slog.InfoContext(ctx, "Entry: HasPendingRequest", "user1_id", user1ID, "user2_id", user2ID)
	var exists bool
	err := repo.database.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM friend_requests WHERE ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) AND status = ?)",
		user1ID, user2ID, user2ID, user1ID, models.FriendRequestPending,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("checking pending request: %w", err)
	}

	slog.InfoContext(ctx, "Exit: HasPendingRequest", "user1_id", user1ID, "user2_id", user2ID, "has_pending", exists)
	return exists, nil
}

func (repo *friendRepository) AcceptFriendRequest(ctx context.Context, requestID string, friendship *models.Friendship) error {
	slog.InfoContext(ctx, "Entry: AcceptFriendRequest", "request_id", requestID, "user_id", friendship.UserID, "friend_id", friendship.FriendID)

	transaction, err := repo.database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer transaction.Rollback()

	// 1. Update request status
	_, err = transaction.ExecContext(ctx, "UPDATE friend_requests SET status = ? WHERE id = ?", models.FriendRequestAccepted, requestID)
	if err != nil {
		return fmt.Errorf("updating friend request status: %w", err)
	}

	// 2. Create friendship (both directions)
	query := "INSERT INTO friendships (user_id, friend_id, created_at) VALUES (?, ?, ?)"

	_, err = transaction.ExecContext(ctx, query, friendship.UserID, friendship.FriendID, friendship.CreatedAt)
	if err != nil {
		return fmt.Errorf("inserting friendship direction 1: %w", err)
	}

	_, err = transaction.ExecContext(ctx, query, friendship.FriendID, friendship.UserID, friendship.CreatedAt)
	if err != nil {
		return fmt.Errorf("inserting friendship direction 2: %w", err)
	}

	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	slog.InfoContext(ctx, "Exit: AcceptFriendRequest", "request_id", requestID)
	return nil
}
