package services

import (
	"context"
	"testing"

	"github.com/igorracki/f1/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockF1DataClient struct {
	mock.Mock
}

func (m *MockF1DataClient) GetScheduleByYear(ctx context.Context, year int) ([]models.RaceWeekend, error) {
	args := m.Called(ctx, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.RaceWeekend), args.Error(1)
}

func (m *MockF1DataClient) GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error) {
	args := m.Called(ctx, year, round, sessionType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SessionResults), args.Error(1)
}

func int64Ptr(v int64) *int64 {
	return &v
}

func TestGetSessionResults_Formatting(t *testing.T) {
	mockClient := new(MockF1DataClient)
	service := NewF1Service(mockClient)
	ctx := context.Background()

	// Scenario: Winner and 2nd Place
	mockResults := &models.SessionResults{
		Year:        2023,
		Round:       1,
		SessionType: models.SessionTypeRaceShort,
		Results: []models.DriverResult{
			{
				Position:    1,
				TotalTimeMS: int64Ptr(5400000), // 1:30:00.000
				GapMS:       int64Ptr(0),
				Status:      "Finished",
			},
			{
				Position:    2,
				TotalTimeMS: nil,             // Should be nil for non-winner per new rules
				GapMS:       int64Ptr(12345), // +12.345s
				Status:      "Finished",
			},
			{
				Position:    3,
				TotalTimeMS: nil,
				GapMS:       int64Ptr(65432), // +1:05.432
				Status:      "Finished",
			},
		},
	}

	mockClient.On("GetSessionResults", ctx, 2023, 1, models.SessionTypeRaceShort).Return(mockResults, nil)

	result, err := service.GetSessionResults(ctx, 2023, 1, models.SessionTypeRaceShort)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify Winner
	winner := result.Results[0]
	assert.Equal(t, "1:30:00.000", winner.TotalTime)
	assert.Equal(t, "+0.000", winner.Gap)

	// Verify 2nd Place
	second := result.Results[1]
	assert.Empty(t, second.TotalTime) // Should be empty string as pointer was nil
	assert.Equal(t, "+12.345", second.Gap)

	// Verify 3rd Place (Gap > 1 minute)
	third := result.Results[2]
	assert.Equal(t, "+1:05.432", third.Gap)
}

func TestGetSessionResults_Qualifying(t *testing.T) {
	mockClient := new(MockF1DataClient)
	service := NewF1Service(mockClient)
	ctx := context.Background()

	mockResults := &models.SessionResults{
		Year:        2023,
		Round:       1,
		SessionType: models.SessionTypeQualifyingShort,
		Results: []models.DriverResult{
			{
				Position: 1,
				Qualifying: &models.QualifyingDetails{
					Q1MS: int64Ptr(90123), // 1:30.123
					Q2MS: int64Ptr(89456), // 1:29.456
					Q3MS: int64Ptr(88789), // 1:28.789
				},
			},
			{
				Position: 15,
				Qualifying: &models.QualifyingDetails{
					Q1MS: int64Ptr(91000), // 1:31.000
					Q2MS: nil,             // Knocked out
					Q3MS: nil,
				},
			},
		},
	}

	mockClient.On("GetSessionResults", ctx, 2023, 1, models.SessionTypeQualifyingShort).Return(mockResults, nil)

	result, err := service.GetSessionResults(ctx, 2023, 1, models.SessionTypeQualifyingShort)
	assert.NoError(t, err)

	// Verify Pole Position
	pole := result.Results[0]
	assert.Equal(t, "1:30.123", pole.Qualifying.Q1)
	assert.Equal(t, "1:29.456", pole.Qualifying.Q2)
	assert.Equal(t, "1:28.789", pole.Qualifying.Q3)

	// Verify Q2 Knockout
	q2out := result.Results[1]
	assert.Equal(t, "1:31.000", q2out.Qualifying.Q1)
	assert.Empty(t, q2out.Qualifying.Q2)
	assert.Empty(t, q2out.Qualifying.Q3)
}

func TestGetSessionResults_NilCheck(t *testing.T) {
	mockClient := new(MockF1DataClient)
	service := NewF1Service(mockClient)
	ctx := context.Background()

	// Scenario: Driver with no time data (e.g., DNF on lap 1)
	mockResults := &models.SessionResults{
		Results: []models.DriverResult{
			{
				Position:    20,
				Status:      "Collision",
				TotalTimeMS: nil,
				GapMS:       nil,
			},
		},
	}

	mockClient.On("GetSessionResults", ctx, 2023, 1, models.SessionTypeRaceShort).Return(mockResults, nil)

	result, err := service.GetSessionResults(ctx, 2023, 1, models.SessionTypeRaceShort)
	assert.NoError(t, err)

	dnf := result.Results[0]
	assert.Equal(t, "Collision", dnf.Gap) // Should fallback to status
	assert.Empty(t, dnf.TotalTime)
}
