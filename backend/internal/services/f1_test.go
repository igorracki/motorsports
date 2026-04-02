package services

import (
	"context"
	"testing"

	"github.com/igorracki/motorsports/backend/internal/models"
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

func (m *MockF1DataClient) GetCircuit(ctx context.Context, year int, round int) (*models.Circuit, error) {
	args := m.Called(ctx, year, round)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Circuit), args.Error(1)
}

func (m *MockF1DataClient) GetDrivers(ctx context.Context, year int, round int) ([]models.DriverInfo, error) {
	args := m.Called(ctx, year, round)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.DriverInfo), args.Error(1)
}

func TestGetScheduleByYear(t *testing.T) {
	mockClient := new(MockF1DataClient)
	baseService := NewF1Service(mockClient)
	service := NewF1CachingService(baseService)
	defer service.Close()
	ctx := context.Background()

	mockSchedule := []models.RaceWeekend{
		{
			Round: 1,
			Name:  "Test GP",
			Sessions: []models.Session{
				{
					Type:        "FP1",
					TimeUTCMS:   500,
					UTCOffsetMS: 1000,
				},
				{
					Type:        "Race",
					TimeUTCMS:   2500,
					UTCOffsetMS: 1000,
				},
			},
		},
	}

	mockClient.On("GetScheduleByYear", ctx, 2024).Return(mockSchedule, nil)

	// When
	result, err := service.GetScheduleByYear(ctx, 2024)

	// Then
	assert.NoError(t, err)
	assert.Len(t, result, 1)

	weekend := result[0]
	// Start local = 500 + 1000 = 1500
	assert.Equal(t, int64(1500), weekend.StartDateLocalMS)
	// End local = 2500 + 1000 = 3500
	assert.Equal(t, int64(3500), weekend.EndDateLocalMS)
	assert.Equal(t, int64(500), weekend.StartDateUTCMS)
	assert.Equal(t, int64(2500), weekend.EndDateUTCMS)

	assert.NotEmpty(t, weekend.StartDateLocal)
	assert.NotEmpty(t, weekend.EndDateLocal)
	assert.NotEmpty(t, weekend.StartDateUTC)
	assert.NotEmpty(t, weekend.EndDateUTC)
}

func TestGetSessionResults_Formatting(t *testing.T) {
	mockClient := new(MockF1DataClient)
	service := NewF1CachingService(NewF1Service(mockClient))
	ctx := context.Background()

	mockSchedule := []models.RaceWeekend{
		{
			Round: 1,
			Sessions: []models.Session{
				{Type: "Race", TimeUTCMS: 1677942400000},
			},
		},
	}
	mockClient.On("GetScheduleByYear", ctx, 2023).Return(mockSchedule, nil)

	mockResults := &models.SessionResults{
		Year:        2023,
		Round:       1,
		SessionType: models.SessionTypeRaceShort,
		Results: []models.DriverResult{
			{
				Position:    1,
				TotalTimeMS: int64Ptr(5400000),
				GapMS:       int64Ptr(0),
				Status:      "Finished",
			},
			{
				Position:    2,
				TotalTimeMS: nil,
				GapMS:       int64Ptr(12345),
				Status:      "Finished",
			},
			{
				Position:    3,
				TotalTimeMS: nil,
				GapMS:       int64Ptr(65432),
				Status:      "Finished",
			},
		},
	}

	mockClient.On("GetSessionResults", ctx, 2023, 1, models.SessionTypeRaceShort).Return(mockResults, nil)

	// When
	result, err := service.GetSessionResults(ctx, 2023, 1, models.SessionTypeRaceShort)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, result)

	winner := result.Results[0]
	assert.Equal(t, "1:30:00.000", winner.TotalTime)
	assert.Equal(t, "+0.000", winner.Gap)

	second := result.Results[1]
	assert.Empty(t, second.TotalTime)
	assert.Equal(t, "+12.345", second.Gap)

	third := result.Results[2]
	assert.Equal(t, "+1:05.432", third.Gap)
}

func TestGetSessionResults_Qualifying(t *testing.T) {
	mockClient := new(MockF1DataClient)
	service := NewF1CachingService(NewF1Service(mockClient))
	ctx := context.Background()

	mockSchedule := []models.RaceWeekend{
		{
			Round: 1,
			Sessions: []models.Session{
				{Type: "Qualifying", TimeUTCMS: 1677942400000},
			},
		},
	}
	mockClient.On("GetScheduleByYear", ctx, 2023).Return(mockSchedule, nil)

	mockResults := &models.SessionResults{
		Year:        2023,
		Round:       1,
		SessionType: models.SessionTypeQualifyingShort,
		Results: []models.DriverResult{
			{
				Position: 1,
				Qualifying: &models.QualifyingDetails{
					Q1MS: int64Ptr(90123),
					Q2MS: int64Ptr(89456),
					Q3MS: int64Ptr(88789),
				},
			},
			{
				Position: 15,
				Qualifying: &models.QualifyingDetails{
					Q1MS: int64Ptr(91000),
					Q2MS: nil,
					Q3MS: nil,
				},
			},
		},
	}

	mockClient.On("GetSessionResults", ctx, 2023, 1, models.SessionTypeQualifyingShort).Return(mockResults, nil)

	// When
	result, err := service.GetSessionResults(ctx, 2023, 1, models.SessionTypeQualifyingShort)

	// Then
	assert.NoError(t, err)

	pole := result.Results[0]
	assert.Equal(t, "1:30.123", pole.Qualifying.Q1)
	assert.Equal(t, "1:29.456", pole.Qualifying.Q2)
	assert.Equal(t, "1:28.789", pole.Qualifying.Q3)

	q2out := result.Results[1]
	assert.Equal(t, "1:31.000", q2out.Qualifying.Q1)
	assert.Empty(t, q2out.Qualifying.Q2)
	assert.Empty(t, q2out.Qualifying.Q3)
}

func TestGetSessionResults_NilCheck(t *testing.T) {
	mockClient := new(MockF1DataClient)
	service := NewF1CachingService(NewF1Service(mockClient))
	ctx := context.Background()

	mockSchedule := []models.RaceWeekend{
		{
			Round: 1,
			Sessions: []models.Session{
				{Type: "Race", TimeUTCMS: 1677942400000},
			},
		},
	}
	mockClient.On("GetScheduleByYear", ctx, 2023).Return(mockSchedule, nil)

	mockResults := &models.SessionResults{
		SessionType: models.SessionTypeRaceShort,
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

	// When
	result, err := service.GetSessionResults(ctx, 2023, 1, models.SessionTypeRaceShort)

	// Then
	assert.NoError(t, err)

	dnf := result.Results[0]
	assert.Equal(t, "Collision", dnf.Gap)
	assert.Empty(t, dnf.TotalTime)
}

func TestGetCircuit(t *testing.T) {
	mockClient := new(MockF1DataClient)
	service := NewF1CachingService(NewF1Service(mockClient))
	ctx := context.Background()

	mockSchedule := []models.RaceWeekend{
		{
			Round:          10,
			StartDateUTCMS: 1677942400000,
			EndDateUTCMS:   1678142400000,
		},
	}
	mockClient.On("GetScheduleByYear", ctx, 2023).Return(mockSchedule, nil)

	mockCircuit := &models.Circuit{
		CircuitName: "Silverstone Circuit",
		Location:    "Silverstone",
		Country:     "UK",
		Latitude:    float64Ptr(52.0786),
		Longitude:   float64Ptr(-1.01694),
		LengthKm:    float64Ptr(5.891),
		Corners:     intPtr(18),
		Layout: []models.CircuitLayoutPoint{
			{X: 100, Y: 200},
			{X: 150, Y: 250},
		},
	}

	mockClient.On("GetCircuit", ctx, 2023, 10).Return(mockCircuit, nil)

	// When
	result, err := service.GetCircuit(ctx, 2023, 10)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Silverstone Circuit", result.CircuitName)
	assert.Equal(t, 52.0786, *result.Latitude)
	assert.Equal(t, -1.01694, *result.Longitude)
	assert.Len(t, result.Layout, 2)
}

func TestGetSessionResults_DerivedGaps(t *testing.T) {
	mockClient := new(MockF1DataClient)
	service := NewF1CachingService(NewF1Service(mockClient))
	ctx := context.Background()

	mockSchedule := []models.RaceWeekend{
		{
			Round: 1,
			Sessions: []models.Session{
				{Type: "Qualifying", TimeUTCMS: 1677942400000},
			},
		},
	}
	mockClient.On("GetScheduleByYear", ctx, 2023).Return(mockSchedule, nil)

	mockResults := &models.SessionResults{
		Year:        2023,
		Round:       1,
		SessionType: models.SessionTypeQualifyingShort,
		Results: []models.DriverResult{
			{
				Position:     1,
				FastestLapMS: int64Ptr(90000), // 1:30.000
			},
			{
				Position:     2,
				FastestLapMS: int64Ptr(90100), // 1:30.100
			},
		},
	}

	mockClient.On("GetSessionResults", ctx, 2023, 1, models.SessionTypeQualifyingShort).Return(mockResults, nil)

	// When
	result, err := service.GetSessionResults(ctx, 2023, 1, models.SessionTypeQualifyingShort)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "-", result.Results[0].Gap)
	assert.Equal(t, "+0.100", result.Results[1].Gap)
}

func TestGetSessionResults_Sprint(t *testing.T) {
	mockClient := new(MockF1DataClient)
	service := NewF1CachingService(NewF1Service(mockClient))
	ctx := context.Background()

	mockSchedule := []models.RaceWeekend{
		{
			Round: 1,
			Sessions: []models.Session{
				{Type: "Sprint", TimeUTCMS: 1677942400000},
			},
		},
	}
	mockClient.On("GetScheduleByYear", ctx, 2023).Return(mockSchedule, nil)

	mockResults := &models.SessionResults{
		Year:        2023,
		Round:       1,
		SessionType: models.SessionTypeSprintShort,
		Results: []models.DriverResult{
			{
				Position:    1,
				TotalTimeMS: int64Ptr(1800000), // 30:00.000
				GapMS:       int64Ptr(0),
				Status:      "Finished",
			},
			{
				Position:    2,
				TotalTimeMS: nil,
				GapMS:       int64Ptr(500), // +0.500
				Status:      "Finished",
			},
		},
	}

	mockClient.On("GetSessionResults", ctx, 2023, 1, models.SessionTypeSprintShort).Return(mockResults, nil)

	// When
	result, err := service.GetSessionResults(ctx, 2023, 1, models.SessionTypeSprintShort)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "30:00.000", result.Results[0].TotalTime)
	assert.Equal(t, "+0.000", result.Results[0].Gap)
	assert.Equal(t, "+0.500", result.Results[1].Gap)
}

func TestGetSessionResults_QualifyingReference(t *testing.T) {
	mockClient := new(MockF1DataClient)
	service := NewF1CachingService(NewF1Service(mockClient))
	ctx := context.Background()

	mockSchedule := []models.RaceWeekend{
		{
			Round: 1,
			Sessions: []models.Session{
				{Type: "Qualifying", TimeUTCMS: 1677942400000},
			},
		},
	}
	mockClient.On("GetScheduleByYear", ctx, 2023).Return(mockSchedule, nil)

	mockResults := &models.SessionResults{
		Year:        2023,
		Round:       1,
		SessionType: models.SessionTypeQualifyingShort,
		Results: []models.DriverResult{
			{
				Position:     1,
				FastestLapMS: int64Ptr(91000), // 1:31.000 (Pole in wet Q3)
			},
			{
				Position:     5,
				FastestLapMS: int64Ptr(90000), // 1:30.000 (Fastest in dry Q2)
			},
		},
	}

	mockClient.On("GetSessionResults", ctx, 2023, 1, models.SessionTypeQualifyingShort).Return(mockResults, nil)

	// When
	result, err := service.GetSessionResults(ctx, 2023, 1, models.SessionTypeQualifyingShort)

	// Then
	assert.NoError(t, err)
	// Reference should be P1 (91000), not the absolute fastest (90000)
	assert.Equal(t, "-", result.Results[0].Gap)
	assert.Equal(t, "-1.000", result.Results[1].Gap) // P5 still in top 10, shows gap
}

func TestGetSessionResults_QualifyingGaps(t *testing.T) {
	mockClient := new(MockF1DataClient)
	service := NewF1CachingService(NewF1Service(mockClient))
	ctx := context.Background()

	mockSchedule := []models.RaceWeekend{
		{
			Round: 1,
			Sessions: []models.Session{
				{Type: "Qualifying", TimeUTCMS: 1677942400000},
			},
		},
	}
	mockClient.On("GetScheduleByYear", ctx, 2023).Return(mockSchedule, nil)

	mockResults := &models.SessionResults{
		Year:        2023,
		Round:       1,
		SessionType: models.SessionTypeQualifyingShort,
		Results: []models.DriverResult{
			{
				Position:     1,
				FastestLapMS: int64Ptr(88789), // 1:28.789
			},
			{
				Position:     2,
				FastestLapMS: int64Ptr(88889), // 1:28.889
			},
			{
				Position:     11,
				FastestLapMS: int64Ptr(89789), // 1:29.789
			},
		},
	}

	mockClient.On("GetSessionResults", ctx, 2023, 1, models.SessionTypeQualifyingShort).Return(mockResults, nil)

	// When
	result, err := service.GetSessionResults(ctx, 2023, 1, models.SessionTypeQualifyingShort)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "-", result.Results[0].Gap)
	assert.Equal(t, "+0.100", result.Results[1].Gap)
	assert.Equal(t, "+1.000", result.Results[2].Gap)
}

func TestGetSessionResults_QualifyingFallback(t *testing.T) {
	mockClient := new(MockF1DataClient)
	service := NewF1CachingService(NewF1Service(mockClient))
	ctx := context.Background()

	mockSchedule := []models.RaceWeekend{
		{
			Round: 1,
			Sessions: []models.Session{
				{Type: "Qualifying", TimeUTCMS: 1677942400000},
			},
		},
	}
	mockClient.On("GetScheduleByYear", ctx, 2023).Return(mockSchedule, nil)

	mockResults := &models.SessionResults{
		Year:        2023,
		Round:       1,
		SessionType: models.SessionTypeQualifyingShort,
		Results: []models.DriverResult{
			{
				Position:     0, // P1 not set correctly
				FastestLapMS: int64Ptr(90000),
			},
			{
				Position:     2,
				FastestLapMS: int64Ptr(90100),
			},
		},
	}

	mockClient.On("GetSessionResults", ctx, 2023, 1, models.SessionTypeQualifyingShort).Return(mockResults, nil)

	// When
	result, err := service.GetSessionResults(ctx, 2023, 1, models.SessionTypeQualifyingShort)

	// Then
	assert.NoError(t, err)
	// Reference should fall back to first driver (index 0)
	assert.Equal(t, "-", result.Results[0].Gap)
	assert.Equal(t, "+0.100", result.Results[1].Gap)
}
