package services

import (
	"time"
)

type PredictionPolicy interface {
	IsLocked(sessionTimeMS int64) bool
	IsScoringFinal(sessionTimeMS int64) bool
	GetRevalidationDeadline(sessionTimeMS int64) time.Time
}

type predictionPolicy struct {
	revalidationWindow time.Duration
	lockThreshold      time.Duration
}

func NewPredictionPolicy() PredictionPolicy {
	return &predictionPolicy{
		revalidationWindow: 48 * time.Hour,
		lockThreshold:      2 * time.Hour,
	}
}

func (p *predictionPolicy) IsLocked(sessionTimeMS int64) bool {
	return time.Now().UTC().UnixMilli() >= sessionTimeMS
}

func (p *predictionPolicy) IsScoringFinal(sessionTimeMS int64) bool {
	deadline := sessionTimeMS + p.revalidationWindow.Milliseconds()
	return time.Now().UTC().UnixMilli() > deadline
}

func (p *predictionPolicy) GetRevalidationDeadline(sessionTimeMS int64) time.Time {
	return time.UnixMilli(sessionTimeMS + p.revalidationWindow.Milliseconds()).UTC()
}

func (p *predictionPolicy) isSessionCompleted(sessionTimeMS int64) bool {
	return time.Now().UTC().UnixMilli() > sessionTimeMS+p.lockThreshold.Milliseconds()
}
