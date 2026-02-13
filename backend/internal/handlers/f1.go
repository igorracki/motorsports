package handlers

import (
	"net/http"
	"strconv"

	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type F1Handler struct {
	service services.F1Service
}

func NewF1Handler(service services.F1Service) *F1Handler {
	return &F1Handler{
		service: service,
	}
}

func (handler *F1Handler) GetRaceWeekends(context echo.Context) error {
	yearParameter := context.Param("year")
	if yearParameter == "" {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_parameter",
			Message: "must provide a year",
		})
	}

	year, err := strconv.Atoi(yearParameter)
	if err != nil {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_parameter",
			Message: "year must be a valid integer",
		})
	}

	if year <= 0 {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_parameter",
			Message: "year must be a positive integer",
		})
	}

	raceWeekends, err := handler.service.GetRaceWeekendsByYear(context.Request().Context(), year)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "failed to fetch race weekends",
		})
	}

	response := models.RaceWeekendsResponse{
		RaceWeekends: raceWeekends,
	}

	return context.JSON(http.StatusOK, response)
}

func (handler *F1Handler) GetSessionResults(context echo.Context) error {
	yearStr := context.Param("year")
	roundStr := context.Param("round")
	sessionType := context.Param("session")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_year"})
	}

	round, err := strconv.Atoi(roundStr)
	if err != nil {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_round"})
	}

	results, err := handler.service.GetSessionResults(context.Request().Context(), year, round, sessionType)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}

	return context.JSON(http.StatusOK, results)
}
