package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MockValidator struct {
	validator *validator.Validate
}

func (mv *MockValidator) Validate(i interface{}) error {
	return mv.validator.Struct(i)
}

func TestF1Handler_Validation(t *testing.T) {
	e := echo.New()
	e.Validator = &MockValidator{validator: validator.New()}
	handler := NewF1Handler(nil)

	t.Run("GetSchedule - Invalid Year", func(tt *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/schedule/1949", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/schedule/:year")
		c.SetParamNames("year")
		c.SetParamValues("1949")

		err := handler.GetSchedule(c)

		assert.Error(tt, err)
		assert.True(tt, errors.Is(err, models.ErrInvalidInput))
	})

	t.Run("GetSessionResults - Invalid Round", func(tt *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/schedule/2024/0/race/results", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/schedule/:year/:round/:session/results")
		c.SetParamNames("year", "round", "session")
		c.SetParamValues("2024", "0", "race")

		err := handler.GetSessionResults(c)

		assert.Error(tt, err)
		assert.True(tt, errors.Is(err, models.ErrInvalidInput))
	})
}
