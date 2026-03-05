package handlers

import (
	"log/slog"
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

func (handler *F1Handler) GetSchedule(context echo.Context) error {
	ctx := context.Request().Context()
	yearParameter := context.Param("year")
	slog.InfoContext(ctx, "Entry: GetSchedule", "year_param", yearParameter)

	if yearParameter == "" {
		slog.WarnContext(ctx, "Missing year parameter")
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_parameter",
			Message: "must provide a year",
		})
	}

	year, err := strconv.Atoi(yearParameter)
	if err != nil {
		slog.WarnContext(ctx, "Invalid year parameter", "year_param", yearParameter)
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_parameter",
			Message: "year must be a valid integer",
		})
	}

	if year < 1950 || year > 2100 {
		slog.WarnContext(ctx, "Year parameter out of bounds", "year", year)
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_parameter",
			Message: "year must be between 1950 and 2100",
		})
	}

	schedule, err := handler.service.GetScheduleByYear(ctx, year)
	if err != nil {
		slog.ErrorContext(ctx, "Service error fetching schedule", "year", year, "error", err)
		return context.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "failed to fetch schedule",
		})
	}

	response := models.ScheduleResponse{
		Schedule: schedule,
	}

	slog.InfoContext(ctx, "Exit: GetSchedule", "year", year, "count", len(schedule))
	return context.JSON(http.StatusOK, response)
}

func (handler *F1Handler) GetSessionResults(context echo.Context) error {
	ctx := context.Request().Context()
	yearStr := context.Param("year")
	roundStr := context.Param("round")
	sessionType := context.Param("session")

	slog.InfoContext(ctx, "Entry: GetSessionResults", "year", yearStr, "round", roundStr, "session", sessionType)

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		slog.WarnContext(ctx, "Invalid year parameter", "year", yearStr)
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_year"})
	}

	if year < 1950 || year > 2100 {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "year_out_of_bounds"})
	}

	round, err := strconv.Atoi(roundStr)
	if err != nil {
		slog.WarnContext(ctx, "Invalid round parameter", "round", roundStr)
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_round"})
	}

	if round < 1 || round > 50 {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "round_out_of_bounds"})
	}

	results, err := handler.service.GetSessionResults(ctx, year, round, sessionType)
	if err != nil {
		slog.ErrorContext(ctx, "Service error fetching session results", "year", year, "round", round, "session", sessionType, "error", err)
		return context.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}

	count := 0
	if results != nil {
		count = len(results.Results)
	}

	slog.InfoContext(ctx, "Exit: GetSessionResults", "year", year, "round", round, "session", sessionType, "count", count)
	return context.JSON(http.StatusOK, results)
}

func (handler *F1Handler) GetCircuit(context echo.Context) error {
	ctx := context.Request().Context()
	yearStr := context.Param("year")
	roundStr := context.Param("round")

	slog.InfoContext(ctx, "Entry: GetCircuit", "year", yearStr, "round", roundStr)

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		slog.WarnContext(ctx, "Invalid year parameter", "year", yearStr)
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_year"})
	}

	if year < 1950 || year > 2100 {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "year_out_of_bounds"})
	}

	round, err := strconv.Atoi(roundStr)
	if err != nil {
		slog.WarnContext(ctx, "Invalid round parameter", "round", roundStr)
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_round"})
	}

	if round < 1 || round > 50 {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "round_out_of_bounds"})
	}

	circuit, err := handler.service.GetCircuit(ctx, year, round)
	if err != nil {
		slog.ErrorContext(ctx, "Service error fetching circuit", "year", year, "round", round, "error", err)
		return context.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}

	if circuit == nil {
		slog.WarnContext(ctx, "Circuit not found", "year", year, "round", round)
		return context.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "circuit not found",
		})
	}

	slog.InfoContext(ctx, "Exit: GetCircuit", "year", year, "round", round, "circuit_name", circuit.CircuitName)
	return context.JSON(http.StatusOK, circuit)
}

func (handler *F1Handler) GetDrivers(context echo.Context) error {
	ctx := context.Request().Context()
	yearStr := context.Param("year")
	roundStr := context.Param("round")

	slog.InfoContext(ctx, "Entry: GetDrivers", "year", yearStr, "round", roundStr)

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		slog.WarnContext(ctx, "Invalid year parameter", "year", yearStr)
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_year"})
	}

	if year < 1950 || year > 2100 {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "year_out_of_bounds"})
	}

	round, err := strconv.Atoi(roundStr)
	if err != nil {
		slog.WarnContext(ctx, "Invalid round parameter", "round", roundStr)
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_round"})
	}

	if round < 1 || round > 50 {
		return context.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "round_out_of_bounds"})
	}

	drivers, err := handler.service.GetDrivers(ctx, year, round)
	if err != nil {
		slog.ErrorContext(ctx, "Service error fetching drivers", "year", year, "round", round, "error", err)
		return context.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}

	slog.InfoContext(ctx, "Exit: GetDrivers", "year", year, "round", round, "count", len(drivers))
	return context.JSON(http.StatusOK, drivers)
}
