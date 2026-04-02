package models

import "time"

type FriendRequestStatus string

const (
	FriendRequestPending  FriendRequestStatus = "pending"
	FriendRequestAccepted FriendRequestStatus = "accepted"
	FriendRequestDenied   FriendRequestStatus = "denied"
)

type FriendRequest struct {
	ID          string              `json:"id"`
	SenderID    string              `json:"sender_id"`
	ReceiverID  string              `json:"receiver_id"`
	Status      FriendRequestStatus `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	SenderEmail string              `json:"sender_email,omitempty"`
	SenderName  string              `json:"sender_name,omitempty"`
}

type Friendship struct {
	UserID    string    `json:"user_id"`
	FriendID  string    `json:"friend_id"`
	CreatedAt time.Time `json:"created_at"`
}

type SendFriendRequestRequest struct {
	Identifier string `json:"identifier" validate:"required"` // Email or UUID
}

type HandleFriendRequestRequest struct {
	Action string `json:"action" validate:"required,oneof=accept deny"`
}

type LeaderboardEntry struct {
	Position    int    `json:"position"`
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Score       int    `json:"score"`
}
