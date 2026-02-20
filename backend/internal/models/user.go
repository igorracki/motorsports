package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Profile struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
}

type UserScore struct {
	UserID    string    `json:"user_id"`
	ScoreType string    `json:"score_type"`
	Season    *int      `json:"season,omitempty"`
	Value     int       `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}
