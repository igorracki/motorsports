package models

import "time"

type User struct {
	ID        string    `json:"id"`
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

type RegisterUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	DisplayName string `json:"display_name" validate:"required"`
}

type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

type AuthResponse struct {
	User    User    `json:"user"`
	Profile Profile `json:"profile"`
}

type UpdateProfileRequest struct {
	DisplayName string `json:"display_name"`
}

type UserProfileResponse struct {
	User    User        `json:"user"`
	Profile Profile     `json:"profile"`
	Scores  []UserScore `json:"scores"`
}
