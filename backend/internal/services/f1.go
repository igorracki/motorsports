package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/igorracki/f1/backend/internal/clients"
	"github.com/igorracki/f1/backend/internal/formatters"
	"github.com/igorracki/f1/backend/internal/models"
)

type F1Service interface {
	GetScheduleByYear(ctx context.Context, year int) ([]models.RaceWeekend, error)
	GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error)
	GetCircuit(ctx context.Context, year int, round int) (*models.Circuit, error)
}

type f1Service struct {
	client clients.F1DataClient
}

func NewF1Service(client clients.F1DataClient) F1Service {
	return &f1Service{
		client: client,
	}
}

func (service *f1Service) GetScheduleByYear(ctx context.Context, year int) ([]models.RaceWeekend, error) {
	slog.InfoContext(ctx, "Entry: GetScheduleByYear", "year", year)
	if year < 1900 || year > 2050 {
		slog.WarnContext(ctx, "Invalid year requested", "year", year)
		return nil, fmt.Errorf("year outside supported Formula 1 range")
	}

	schedule, err := service.client.GetScheduleByYear(ctx, year)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch schedule", "error", err)
		return nil, fmt.Errorf("failed to fetch schedule from external API: %w", err)
	}

	for weekendIndex := range schedule {
		schedule[weekendIndex].StartDate = formatters.FormatTimestamp(schedule[weekendIndex].StartDateMS)
		for sessionIndex := range schedule[weekendIndex].Sessions {
			schedule[weekendIndex].Sessions[sessionIndex].TimeLocal = formatters.FormatTimestamp(schedule[weekendIndex].Sessions[sessionIndex].TimeLocalMS)
			schedule[weekendIndex].Sessions[sessionIndex].TimeUTC = formatters.FormatTimestamp(schedule[weekendIndex].Sessions[sessionIndex].TimeUTCMS)
		}
	}

	slog.InfoContext(ctx, "Exit: GetScheduleByYear", "year", year, "count", len(schedule))
	return schedule, nil
}

func (service *f1Service) GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error) {
	slog.InfoContext(ctx, "Entry: GetSessionResults", "year", year, "round", round, "sessionType", sessionType)
	results, err := service.client.GetSessionResults(ctx, year, round, sessionType)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch session results", "error", err)
		return nil, fmt.Errorf("failed to fetch results for session %s in round %d (%d): %w", sessionType, round, year, err)
	}

	if results == nil || len(results.Results) == 0 {
		slog.WarnContext(ctx, "No results found", "year", year, "round", round)
		slog.InfoContext(ctx, "Exit: GetSessionResults", "year", year, "round", round, "sessionType", sessionType, "count", 0)
		if results == nil {
			return &models.SessionResults{
				Year:        year,
				Round:       round,
				SessionType: sessionType,
				Results:     []models.DriverResult{},
			}, nil
		}
		return results, nil
	}

	for resultIndex := range results.Results {
		result := &results.Results[resultIndex]

		if result.TotalTimeMS != nil {
			result.TotalTime = formatters.FormatDuration(*result.TotalTimeMS, false)
		}

		if result.FastestLapMS != nil {
			result.FastestLap = formatters.FormatDuration(*result.FastestLapMS, false)
		}

		if sessionType == models.SessionTypeRaceShort || sessionType == models.SessionTypeRace {
			if result.GapMS != nil {
				result.Gap = formatters.FormatDuration(*result.GapMS, true)
			} else {
				result.Gap = result.Status
			}
		}

		if result.Race != nil {
			result.Race.PositionsChange = result.Race.GridPosition - result.Position
		}

		if result.Qualifying != nil {
			if result.Qualifying.Q1MS != nil {
				result.Qualifying.Q1 = formatters.FormatDuration(*result.Qualifying.Q1MS, false)
			}
			if result.Qualifying.Q2MS != nil {
				result.Qualifying.Q2 = formatters.FormatDuration(*result.Qualifying.Q2MS, false)
			}
			if result.Qualifying.Q3MS != nil {
				result.Qualifying.Q3 = formatters.FormatDuration(*result.Qualifying.Q3MS, false)
			}
		}
	}

	slog.InfoContext(ctx, "Exit: GetSessionResults", "year", year, "round", round, "sessionType", sessionType, "count", len(results.Results))
	return results, nil
}

func (service *f1Service) GetCircuit(ctx context.Context, year int, round int) (*models.Circuit, error) {
	slog.InfoContext(ctx, "Entry: GetCircuit", "year", year, "round", round)

	circuit, err := service.client.GetCircuit(ctx, year, round)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch circuit", "error", err)
		return nil, fmt.Errorf("failed to fetch circuit for round %d (%d): %w", round, year, err)
	}

	if circuit != nil {
		circuit.EventDate = formatters.FormatTimestamp(circuit.EventDateMS)
		roundCircuitMetrics(circuit)
		transformLayout(circuit)
	}

	slog.InfoContext(ctx, "Exit: GetCircuit", "year", year, "round", round, "circuit_name", circuit.CircuitName)
	return circuit, nil
}
