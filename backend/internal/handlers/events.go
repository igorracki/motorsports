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

func (handler *F1Handler) GetEvents(context echo.Context) error {
	yearParameter := context.QueryParam("year")
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

	events, err := handler.service.GetEventsByYear(context.Request().Context(), year)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "failed to fetch events",
		})
	}

	response := models.EventsResponse{
		Events: events,
	}

	return context.JSON(http.StatusOK, response)
}
