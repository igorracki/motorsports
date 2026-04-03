package services

import (
	"time"

	"github.com/igorracki/motorsports/backend/internal/models"
)

const (
	DefaultRevalidationWindow = 48 * time.Hour
	DefaultLockThreshold      = 0 * time.Second // Lock exactly at start time
	DefaultPreSessionBuffer   = 15 * time.Minute
	DefaultSessionDuration    = 2 * time.Hour
)

type PredictionPolicy interface {
	IsLocked(sessionTimeMS int64) bool
	IsLive(sessionTimeMS int64) bool
	IsScoringFinal(sessionTimeMS int64) bool
	IsSessionCompleted(sessionTimeMS int64) bool
	GetRevalidationDeadline(sessionTimeMS int64) time.Time
	GetConfig() models.PredictionPolicyConfig
}

type predictionPolicy struct {
	revalidationWindow time.Duration
	lockThreshold      time.Duration
	preSessionBuffer   time.Duration
	sessionDuration    time.Duration
}

func NewPredictionPolicy() PredictionPolicy {
	return &predictionPolicy{
		revalidationWindow: DefaultRevalidationWindow,
		lockThreshold:      DefaultLockThreshold,
		preSessionBuffer:   DefaultPreSessionBuffer,
		sessionDuration:    DefaultSessionDuration,
	}
}

func (policy *predictionPolicy) IsLocked(sessionTimeMS int64) bool {
	return time.Now().UTC().UnixMilli() >= sessionTimeMS+policy.lockThreshold.Milliseconds()
}

func (policy *predictionPolicy) IsLive(sessionTimeMS int64) bool {
	now := time.Now().UTC().UnixMilli()
	return now >= sessionTimeMS-policy.preSessionBuffer.Milliseconds() &&
		now <= sessionTimeMS+policy.sessionDuration.Milliseconds()
}

func (policy *predictionPolicy) IsScoringFinal(sessionTimeMS int64) bool {
	deadline := sessionTimeMS + policy.revalidationWindow.Milliseconds()
	return time.Now().UTC().UnixMilli() > deadline
}

func (policy *predictionPolicy) IsSessionCompleted(sessionTimeMS int64) bool {
	return time.Now().UTC().UnixMilli() > sessionTimeMS+policy.sessionDuration.Milliseconds()
}

func (policy *predictionPolicy) GetRevalidationDeadline(sessionTimeMS int64) time.Time {
	return time.UnixMilli(sessionTimeMS + policy.revalidationWindow.Milliseconds()).UTC()
}

func (policy *predictionPolicy) GetConfig() models.PredictionPolicyConfig {
	return models.PredictionPolicyConfig{
		LockThresholdMS:      policy.lockThreshold.Milliseconds(),
		PreSessionBufferMS:   policy.preSessionBuffer.Milliseconds(),
		SessionDurationMS:    policy.sessionDuration.Milliseconds(),
		RevalidationWindowMS: policy.revalidationWindow.Milliseconds(),
	}
}
