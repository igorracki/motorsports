package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/igorracki/motorsports/backend/internal/api"
	"github.com/igorracki/motorsports/backend/internal/clients"
	"github.com/igorracki/motorsports/backend/internal/handlers"
	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockF1DataClient struct {
	scheduleResponse []models.RaceWeekend
	resultsResponse  *models.SessionResults
	driversResponse  []models.DriverInfo
	err              error
}

func (mock *mockF1DataClient) GetScheduleByYear(ctx context.Context, year int) ([]models.RaceWeekend, error) {
	return mock.scheduleResponse, mock.err
}

func (mock *mockF1DataClient) GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error) {
	return mock.resultsResponse, mock.err
}

func (mock *mockF1DataClient) GetCircuit(ctx context.Context, year int, round int) (*models.Circuit, error) {
	return nil, mock.err
}

func (mock *mockF1DataClient) GetDrivers(ctx context.Context, year int, round int) ([]models.DriverInfo, error) {
	return mock.driversResponse, mock.err
}

func setupTestServer(client clients.F1DataClient) *echo.Echo {
	e := echo.New()
	e.Validator = api.NewCustomValidator()
	e.HTTPErrorHandler = api.HTTPErrorHandler

	policy := services.NewPredictionPolicy()
	baseService := services.NewF1Service(client, policy)
	f1DataService := services.NewF1CachingService(baseService)
	f1DataHandler := handlers.NewF1Handler(f1DataService)

	apiGroup := e.Group("/api")
	apiGroup.GET("/schedule/:year", f1DataHandler.GetSchedule)
	apiGroup.GET("/schedule/:year/:round/:session/results", f1DataHandler.GetSessionResults)
	apiGroup.GET("/schedule/:year/:round/circuit", f1DataHandler.GetCircuit)

	return e
}

func TestGetSchedule_Success(t *testing.T) {
	clientMock := &mockF1DataClient{
		scheduleResponse: []models.RaceWeekend{
			{
				Round:    1,
				FullName: "FORMULA 1 QATAR AIRWAYS AUSTRALIAN GRAND PRIX 2026",
				Name:     "Australian Grand Prix",
				Location: "Melbourne",
				Country:  "Australia",
				Sessions: []models.Session{
					{Type: "Race", TimeUTCMS: 1772942400000, UTCOffsetMS: 39600000},
				},
			},
		},
		err: nil,
	}

	server := setupTestServer(clientMock)

	request := httptest.NewRequest(http.MethodGet, "/api/schedule/2026", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response models.ScheduleResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Len(t, response.Schedule, 1)
	assert.Equal(t, 1, response.Schedule[0].Round)
	assert.Equal(t, "Australian Grand Prix", response.Schedule[0].Name)
	assert.Equal(t, "Melbourne", response.Schedule[0].Location)
	assert.Equal(t, "Australia", response.Schedule[0].Country)
	assert.Len(t, response.Schedule[0].Sessions, 1)
	assert.Equal(t, "Race", response.Schedule[0].Sessions[0].Type)
}

func TestGetSchedule_InvalidYear(t *testing.T) {
	server := setupTestServer(&mockF1DataClient{})

	request := httptest.NewRequest(http.MethodGet, "/api/schedule/invalid", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid_input", response.Error)
}

func TestGetSchedule_NegativeYear(t *testing.T) {
	server := setupTestServer(&mockF1DataClient{})

	request := httptest.NewRequest(http.MethodGet, "/api/schedule/-1", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetSchedule_ExternalAPIError(t *testing.T) {
	clientMock := &mockF1DataClient{
		scheduleResponse: nil,
		err:              assert.AnError,
	}

	server := setupTestServer(clientMock)

	request := httptest.NewRequest(http.MethodGet, "/api/schedule/2026", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "internal_error", response.Error)
}

func TestGetSchedule_MultipleEvents(t *testing.T) {
	clientMock := &mockF1DataClient{
		scheduleResponse: []models.RaceWeekend{
			{
				Round:    1,
				Name:     "Australian Grand Prix",
				Location: "Melbourne",
				Country:  "Australia",
				Sessions: []models.Session{},
			},
			{
				Round:    2,
				Name:     "Chinese Grand Prix",
				Location: "Shanghai",
				Country:  "China",
				Sessions: []models.Session{},
			},
		},
	}

	server := setupTestServer(clientMock)

	request := httptest.NewRequest(http.MethodGet, "/api/schedule/2026", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response models.ScheduleResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Len(t, response.Schedule, 2)
	assert.Equal(t, "Australian Grand Prix", response.Schedule[0].Name)
	assert.Equal(t, "Chinese Grand Prix", response.Schedule[1].Name)
}

func TestGetSessionResults_Success(t *testing.T) {
	mockResults := &models.SessionResults{
		Year:        2023,
		Round:       1,
		SessionType: "Race",
		Results: []models.DriverResult{
			{
				Position: 1,
				Driver: models.DriverInfo{
					ID:     "VER",
					Number: "1",
				},
				Status: "Finished",
			},
		},
	}

	clientMock := &mockF1DataClient{
		scheduleResponse: []models.RaceWeekend{
			{Round: 1, Name: "Test GP"},
		},
		resultsResponse: mockResults,
		err:             nil,
	}

	server := setupTestServer(clientMock)

	request := httptest.NewRequest(http.MethodGet, "/api/schedule/2023/1/Race/results", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response models.SessionResults
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, 2023, response.Year)
	assert.Equal(t, 1, response.Round)
	assert.Equal(t, "Race", response.SessionType)
	assert.Len(t, response.Results, 1)
	assert.Equal(t, 1, response.Results[0].Position)
	assert.Equal(t, "VER", response.Results[0].Driver.ID)
	assert.Equal(t, "1", response.Results[0].Driver.Number)
}
