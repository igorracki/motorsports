package main

import (
	"context"
	"encoding/json"
	"github.com/igorracki/f1/backend/internal/clients"
	"github.com/igorracki/f1/backend/internal/handlers"
	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockF1DataClient struct {
	response []models.Event
	err      error
}

func (mock *mockF1DataClient) GetEventsByYear(ctx context.Context, year int) ([]models.Event, error) {
	return mock.response, mock.err
}

func setupTestServer(client clients.F1DataClient) *echo.Echo {
	server := echo.New()

	eventsService := services.NewF1Service(client)
	eventsHandler := handlers.NewF1Handler(eventsService)

	api := server.Group("/api")
	api.GET("/events", eventsHandler.GetEvents)

	return server
}

func TestGetEvents_Success(t *testing.T) {
	clientMock := &mockF1DataClient{
		response: []models.Event{
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

	request := httptest.NewRequest(http.MethodGet, "/api/events?year=2026", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response models.EventsResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Len(t, response.Events, 1)
	assert.Equal(t, 1, response.Events[0].Round)
	assert.Equal(t, "Australian Grand Prix", response.Events[0].Name)
	assert.Equal(t, "Melbourne", response.Events[0].Location)
	assert.Equal(t, "Australia", response.Events[0].Country)
	assert.Len(t, response.Events[0].Sessions, 1)
	assert.Equal(t, "Race", response.Events[0].Sessions[0].Type)
}
func TestGetEvents_MissingYear(t *testing.T) {
	server := setupTestServer(&mockF1DataClient{})

	request := httptest.NewRequest(http.MethodGet, "/api/events", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "missing_parameter", response.Error)
	assert.Equal(t, "must provide a year", response.Message)
}
func TestGetEvents_InvalidYear(t *testing.T) {
	server := setupTestServer(&mockF1DataClient{})

	request := httptest.NewRequest(http.MethodGet, "/api/events?year=invalid", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid_parameter", response.Error)
}
func TestGetEvents_NegativeYear(t *testing.T) {
	server := setupTestServer(&mockF1DataClient{})

	request := httptest.NewRequest(http.MethodGet, "/api/events?year=-1", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
func TestGetEvents_ExternalAPIError(t *testing.T) {
	clientMock := &mockF1DataClient{
		response: nil,
		err:      assert.AnError,
	}

	server := setupTestServer(clientMock)

	request := httptest.NewRequest(http.MethodGet, "/api/events?year=2026", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "internal_error", response.Error)
}
func TestGetEvents_MultipleEvents(t *testing.T) {
	clientMock := &mockF1DataClient{
		response: []models.Event{
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

	request := httptest.NewRequest(http.MethodGet, "/api/events?year=2026", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response models.EventsResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Len(t, response.Events, 2)
	assert.Equal(t, "Australian Grand Prix", response.Events[0].Name)
	assert.Equal(t, "Chinese Grand Prix", response.Events[1].Name)
}
