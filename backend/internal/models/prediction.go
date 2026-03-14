package models

import "time"

type Prediction struct {
	ID              string            `json:"id"`
	UserID          string            `json:"user_id"`
	Year            int               `json:"year"`
	Round           int               `json:"round"`
	SessionType     string            `json:"session_type"`
	Score           *int              `json:"score,omitempty"`
	RevalidateUntil *time.Time        `json:"revalidate_until,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Entries         []PredictionEntry `json:"entries"`
}

type PredictionEntry struct {
	PredictionID string `json:"prediction_id"`
	Position     int    `json:"position"`
	DriverID     string `json:"driver_id"`
	Correct      bool   `json:"correct"`
}
