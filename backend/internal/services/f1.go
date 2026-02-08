package services

import (
	"context"
	"fmt"

	"github.com/igorracki/f1/backend/internal/clients"
	"github.com/igorracki/f1/backend/internal/models"
)

type F1Service interface {
	GetRaceWeekendsByYear(ctx context.Context, year int) ([]models.RaceWeekend, error)
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

func (service *f1Service) GetRaceWeekendsByYear(ctx context.Context, year int) ([]models.RaceWeekend, error) {
	if year < 1900 || year > 2050 {
		return nil, fmt.Errorf("year outside supported Formula 1 range")
	}

	raceWeekends, err := service.client.GetRaceWeekendsByYear(ctx, year)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch race weekends from external API: %w", err)
	}

	return raceWeekends, nil
}

func (service *f1Service) GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error) {
	results, err := service.client.GetSessionResults(ctx, year, round, sessionType)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch results: %w", err)
	}

	if results == nil || len(results.Results) == 0 {
		return results, nil
	}

	for i := range results.Results {
		res := &results.Results[i]

		if res.TotalTimeMS > 0 {
			res.TotalTime = formatDuration(res.TotalTimeMS, false)
		}

		if res.FastestLapMS > 0 {
			res.FastestLap = formatDuration(res.FastestLapMS, false)
		}

		if sessionType == "R" || sessionType == "Race" {
			if res.Position == 1 {
				res.Gap = "+0.000"
			} else if res.GapMS > 0 {
				res.Gap = formatDuration(res.GapMS, true)
			} else {
				res.Gap = res.Status
			}
		}

		if res.Race != nil {
			res.Race.PositionsChange = res.Race.GridPosition - res.Position
		}

		if res.Qualifying != nil {
			if res.Qualifying.Q1MS > 0 {
				res.Qualifying.Q1 = formatDuration(res.Qualifying.Q1MS, false)
			}
			if res.Qualifying.Q2MS > 0 {
				res.Qualifying.Q2 = formatDuration(res.Qualifying.Q2MS, false)
			}
			if res.Qualifying.Q3MS > 0 {
				res.Qualifying.Q3 = formatDuration(res.Qualifying.Q3MS, false)
			}
		}
	}

	return results, nil
}

func formatDuration(ms int64, isGap bool) string {
	prefix := ""
	if isGap {
		prefix = "+"
	}

	secondsTotal := float64(ms) / 1000.0

	// For very small gaps
	if secondsTotal < 60 && isGap {
		return fmt.Sprintf("%s%.3f", prefix, secondsTotal)
	}

	hours := int(secondsTotal) / 3600
	minutes := (int(secondsTotal) % 3600) / 60
	seconds := secondsTotal - float64(hours*3600+minutes*60)

	if hours > 0 {
		return fmt.Sprintf("%s%d:%02d:%06.3f", prefix, hours, minutes, seconds)
	}

	return fmt.Sprintf("%s%d:%06.3f", prefix, minutes, seconds)
}
