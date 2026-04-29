package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPredictionService struct {
	mock.Mock
	services.PredictionService
}

func (m *mockPredictionService) SubmitPrediction(ctx context.Context, p *models.Prediction) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

type mockScoringService struct {
	mock.Mock
	services.ScoringService
}

type dummyValidator struct{}

func (v *dummyValidator) Validate(i interface{}) error { return nil }

func TestPredictionHandler_SubmitValidation(t *testing.T) {
	e := echo.New()
	e.Validator = &dummyValidator{}

	configService := services.NewConfigService()

	serviceMock := &mockPredictionService{}
	scoringMock := &mockScoringService{}
	handler := NewPredictionHandler(serviceMock, scoringMock, configService)

	tests := []struct {
		name       string
		prediction models.Prediction
		wantStatus int
		wantErr    string
	}{
		{
			name: "Valid Prediction",
			prediction: models.Prediction{
				Year: 2024, Round: 1, SessionType: "Race",
				Entries: []models.PredictionEntry{
					{Position: 1, DriverID: "VER"},
					{Position: 2, DriverID: "PER"},
					{Position: 3, DriverID: "ALO"},
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Too few entries",
			prediction: models.Prediction{
				Year: 2024, Round: 1, SessionType: "Race",
				Entries: []models.PredictionEntry{
					{Position: 1, DriverID: "VER"},
					{Position: 2, DriverID: "PER"},
				},
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    "between 3 and 22 entries",
		},
		{
			name: "Duplicate driver",
			prediction: models.Prediction{
				Year: 2024, Round: 1, SessionType: "Race",
				Entries: []models.PredictionEntry{
					{Position: 1, DriverID: "VER"},
					{Position: 2, DriverID: "VER"},
					{Position: 3, DriverID: "ALO"},
				},
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    "duplicate driver VER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.prediction)
			req := httptest.NewRequest(http.MethodPost, "/users/123/predictions", strings.NewReader(string(body)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues("123")

			if tt.wantStatus == http.StatusCreated {
				serviceMock.On("SubmitPrediction", mock.Anything, mock.Anything).Return(nil).Once()
			}

			err := handler.SubmitPrediction(c)

			if tt.wantStatus == http.StatusBadRequest {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, rec.Code)
			}
		})
	}
}
