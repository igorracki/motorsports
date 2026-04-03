package models

import "time"

type Prediction struct {
	ID              string            `json:"id"`
	UserID          string            `json:"user_id"`
	Year            int               `json:"year" param:"year" validate:"required,min=1950,max=2100"`
	Round           int               `json:"round" param:"round" validate:"required,min=1,max=50"`
	SessionType     string            `json:"session_type" param:"session" validate:"required"`
	Score           *int              `json:"score,omitempty"`
	RevalidateUntil *time.Time        `json:"revalidate_until,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Entries         []PredictionEntry `json:"entries" validate:"required,min=3,max=22,dive"`
}

type PredictionEntry struct {
	PredictionID string `json:"prediction_id"`
	Position     int    `json:"position" validate:"required,min=1,max=22"`
	DriverID     string `json:"driver_id" validate:"required"`
	Correct      bool   `json:"correct"`
}

type PredictionPolicyConfig struct {
	LockThresholdMS      int64 `json:"lock_threshold_ms"`
	PreSessionBufferMS   int64 `json:"pre_session_buffer_ms"`
	SessionDurationMS    int64 `json:"session_duration_ms"`
	RevalidationWindowMS int64 `json:"revalidation_window_ms"`
}
