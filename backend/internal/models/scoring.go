package models

type PositionPoints struct {
	Position int `json:"position"`
	Points   int `json:"points"`
}

type SessionScoringRules struct {
	SessionType string           `json:"session_type"`
	Rules       []PositionPoints `json:"rules"`
}
