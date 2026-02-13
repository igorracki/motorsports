package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/igorracki/f1/backend/internal/clients"
	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/utils"
)

type F1Service interface {
	GetScheduleByYear(ctx context.Context, year int) ([]models.RaceWeekend, error)
	GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error)
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
	slog.InfoContext(ctx, "Processing schedule request", "year", year)
	if year < 1900 || year > 2050 {
		slog.WarnContext(ctx, "Invalid year requested", "year", year)
		return nil, fmt.Errorf("year outside supported Formula 1 range")
	}

	schedule, err := service.client.GetScheduleByYear(ctx, year)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch schedule", "error", err)
		return nil, fmt.Errorf("failed to fetch schedule from external API: %w", err)
	}

	return schedule, nil
}

func (service *f1Service) GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error) {
	slog.InfoContext(ctx, "Processing session results request", "year", year, "round", round, "sessionType", sessionType)
	results, err := service.client.GetSessionResults(ctx, year, round, sessionType)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch session results", "error", err)
		return nil, fmt.Errorf("failed to fetch results: %w", err)
	}

	if results == nil || len(results.Results) == 0 {
		slog.WarnContext(ctx, "No results found", "year", year, "round", round)
		return results, nil
	}

	slog.InfoContext(ctx, "Formatting driver results", "count", len(results.Results))
	for i := range results.Results {
		result := &results.Results[i]

		if result.TotalTimeMS != nil {
			result.TotalTime = utils.FormatDuration(*result.TotalTimeMS, false)
		}

		if result.FastestLapMS != nil {
			result.FastestLap = utils.FormatDuration(*result.FastestLapMS, false)
		}

		if sessionType == models.SessionTypeRaceShort || sessionType == models.SessionTypeRace {
			if result.GapMS != nil {
				result.Gap = utils.FormatDuration(*result.GapMS, true)
			} else {
				result.Gap = result.Status
			}
		}

		if result.Race != nil {
			result.Race.PositionsChange = result.Race.GridPosition - result.Position
		}

		if result.Qualifying != nil {
			if result.Qualifying.Q1MS != nil {
				result.Qualifying.Q1 = utils.FormatDuration(*result.Qualifying.Q1MS, false)
			}
			if result.Qualifying.Q2MS != nil {
				result.Qualifying.Q2 = utils.FormatDuration(*result.Qualifying.Q2MS, false)
			}
			if result.Qualifying.Q3MS != nil {
				result.Qualifying.Q3 = utils.FormatDuration(*result.Qualifying.Q3MS, false)
			}
		}
	}

	return results, nil
}
