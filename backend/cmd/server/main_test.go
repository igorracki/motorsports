package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/igorracki/f1/backend/internal/clients"
	"github.com/igorracki/f1/backend/internal/handlers"
	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockF1DataClient struct {
	raceWeekendsResponse        []models.RaceWeekend
	resultsResponse *models.SessionResults
	err             error
}

func (mock *mockF1DataClient) GetRaceWeekendsByYear(ctx context.Context, year int) ([]models.RaceWeekend, error) {
	return mock.raceWeekendsResponse, mock.err
}

func (mock *mockF1DataClient) GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error) {
	return mock.resultsResponse, mock.err
}

func setupTestServer(client clients.F1DataClient) *echo.Echo {
	server := echo.New()

	eventsService := services.NewF1Service(client)
	eventsHandler := handlers.NewF1Handler(eventsService)

	api := server.Group("/api")
	api.GET("/race-weekends/:year", eventsHandler.GetRaceWeekends)

	return server
}

func TestGetRaceWeekends_Success(t *testing.T) {
	clientMock := &mockF1DataClient{
		raceWeekendsResponse: []models.RaceWeekend{
			{
				Round:     1,
				FullName:  "FORMULA 1 QATAR AIRWAYS AUSTRALIAN GRAND PRIX 2026",
				Name:      "Australian Grand Prix",
				Location:  "Melbourne",
				Country:   "Australia",
				StartDate: "2026-03-08T00:00:00",
				Sessions: []models.Session{
					{Type: "Race", TimeLocal: "2026-03-08T15:00:00+11:00", TimeUTC: "2026-03-08T04:00:00"},
				},
			},
		},
		err: nil,
	}

	server := setupTestServer(clientMock)

	request := httptest.NewRequest(http.MethodGet, "/api/race-weekends/2026", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response models.RaceWeekendsResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Len(t, response.RaceWeekends, 1)
	assert.Equal(t, 1, response.RaceWeekends[0].Round)
	assert.Equal(t, "Australian Grand Prix", response.RaceWeekends[0].Name)
	assert.Equal(t, "Melbourne", response.RaceWeekends[0].Location)
	assert.Equal(t, "Australia", response.RaceWeekends[0].Country)
	assert.Len(t, response.RaceWeekends[0].Sessions, 1)
	assert.Equal(t, "Race", response.RaceWeekends[0].Sessions[0].Type)
}
func TestGetRaceWeekends_InvalidYear(t *testing.T) {
	server := setupTestServer(&mockF1DataClient{})

	request := httptest.NewRequest(http.MethodGet, "/api/race-weekends/invalid", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid_parameter", response.Error)
}
func TestGetRaceWeekends_NegativeYear(t *testing.T) {
	server := setupTestServer(&mockF1DataClient{})

	request := httptest.NewRequest(http.MethodGet, "/api/race-weekends/-1", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
func TestGetRaceWeekends_ExternalAPIError(t *testing.T) {
	clientMock := &mockF1DataClient{
		raceWeekendsResponse: nil,
		err:      assert.AnError,
	}

	server := setupTestServer(clientMock)

	request := httptest.NewRequest(http.MethodGet, "/api/race-weekends/2026", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "internal_error", response.Error)
}
func TestGetRaceWeekends_MultipleEvents(t *testing.T) {
	clientMock := &mockF1DataClient{
		raceWeekendsResponse: []models.RaceWeekend{
			{
				Round:     1,
				Name:      "Australian Grand Prix",
				Location:  "Melbourne",
				Country:   "Australia",
				StartDate: "2026-03-08T00:00:00",
				Sessions:  []models.Session{},
			},
			{
				Round:     2,
				Name:      "Chinese Grand Prix",
				Location:  "Shanghai",
				Country:   "China",
				StartDate: "2026-03-15T00:00:00",
				Sessions:  []models.Session{},
			},
		},
	}

	server := setupTestServer(clientMock)

	request := httptest.NewRequest(http.MethodGet, "/api/race-weekends/2026", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response models.RaceWeekendsResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Len(t, response.RaceWeekends, 2)
	assert.Equal(t, "Australian Grand Prix", response.RaceWeekends[0].Name)
	assert.Equal(t, "Chinese Grand Prix", response.RaceWeekends[1].Name)
}
