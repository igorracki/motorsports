package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestF1Handler_Validation(t *testing.T) {
	e := echo.New()
	handler := NewF1Handler(nil)

	t.Run("GetSchedule - Invalid Year", func(tt *testing.T) {
		// Given
		req := httptest.NewRequest(http.MethodGet, "/api/schedule/1949", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/schedule/:year")
		c.SetParamNames("year")
		c.SetParamValues("1949")

		// When
		err := handler.GetSchedule(c)

		// Then
		if assert.NoError(tt, err) {
			assert.Equal(tt, http.StatusBadRequest, rec.Code)
			assert.Contains(tt, rec.Body.String(), "year must be between 1950 and 2100")
		}
	})

	t.Run("GetSessionResults - Invalid Round", func(tt *testing.T) {
		// Given
		req := httptest.NewRequest(http.MethodGet, "/api/schedule/2024/0/race/results", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/schedule/:year/:round/:session/results")
		c.SetParamNames("year", "round", "session")
		c.SetParamValues("2024", "0", "race")

		// When
		err := handler.GetSessionResults(c)

		// Then
		if assert.NoError(tt, err) {
			assert.Equal(tt, http.StatusBadRequest, rec.Code)
			assert.Contains(tt, rec.Body.String(), "round_out_of_bounds")
		}
	})
}
