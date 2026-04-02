package handlers

import (
	"net/http"

	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/igorracki/motorsports/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type F1Handler struct {
	service services.F1Service
}

type YearParam struct {
	Year int `param:"year" validate:"required,min=1950,max=2100"`
}

type SessionResultsParam struct {
	Year    int    `param:"year" validate:"required,min=1950,max=2100"`
	Round   int    `param:"round" validate:"required,min=1,max=50"`
	Session string `param:"session" validate:"required"`
}

type YearRoundParam struct {
	Year  int `param:"year" validate:"required,min=1950,max=2100"`
	Round int `param:"round" validate:"required,min=1,max=50"`
}

func NewF1Handler(service services.F1Service) *F1Handler {
	return &F1Handler{
		service: service,
	}
}

func (handler *F1Handler) GetSchedule(context echo.Context) error {
	ctx := context.Request().Context()

	var params YearParam
	if err := context.Bind(&params); err != nil {
		return models.ErrInvalidInput
	}
	if err := context.Validate(&params); err != nil {
		return models.ErrInvalidInput
	}

	schedule, err := handler.service.GetScheduleByYear(ctx, params.Year)
	if err != nil {
		return err
	}

	if schedule == nil {
		schedule = []models.RaceWeekend{}
	}

	return context.JSON(http.StatusOK, models.ScheduleResponse{
		Schedule: schedule,
	})
}

func (handler *F1Handler) GetSessionResults(context echo.Context) error {
	ctx := context.Request().Context()

	var params SessionResultsParam
	if err := context.Bind(&params); err != nil {
		return models.ErrInvalidInput
	}
	if err := context.Validate(&params); err != nil {
		return models.ErrInvalidInput
	}

	results, err := handler.service.GetSessionResults(ctx, params.Year, params.Round, params.Session)
	if err != nil {
		return err
	}

	if results == nil {
		return models.ErrNotFound
	}

	if results.Results == nil {
		results.Results = []models.DriverResult{}
	}

	return context.JSON(http.StatusOK, results)
}

func (handler *F1Handler) GetCircuit(context echo.Context) error {
	ctx := context.Request().Context()

	var params YearRoundParam
	if err := context.Bind(&params); err != nil {
		return models.ErrInvalidInput
	}
	if err := context.Validate(&params); err != nil {
		return models.ErrInvalidInput
	}

	circuit, err := handler.service.GetCircuit(ctx, params.Year, params.Round)
	if err != nil {
		return err
	}

	if circuit == nil {
		return models.ErrNotFound
	}

	return context.JSON(http.StatusOK, circuit)
}

func (handler *F1Handler) GetDrivers(context echo.Context) error {
	ctx := context.Request().Context()

	var params YearRoundParam
	if err := context.Bind(&params); err != nil {
		return models.ErrInvalidInput
	}
	if err := context.Validate(&params); err != nil {
		return models.ErrInvalidInput
	}

	drivers, err := handler.service.GetDrivers(ctx, params.Year, params.Round)
	if err != nil {
		return err
	}

	if drivers == nil {
		drivers = []models.DriverInfo{}
	}

	return context.JSON(http.StatusOK, drivers)
}
