package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/igorracki/f1/backend/internal/cache"
	"github.com/igorracki/f1/backend/internal/clients"
	"github.com/igorracki/f1/backend/internal/formatters"
	"github.com/igorracki/f1/backend/internal/models"
)

type F1Service interface {
	GetScheduleByYear(ctx context.Context, year int) ([]models.RaceWeekend, error)
	GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error)
	GetCircuit(ctx context.Context, year int, round int) (*models.Circuit, error)
	GetDrivers(ctx context.Context, year int, round int) ([]models.DriverInfo, error)
}

type f1Service struct {
	client clients.F1DataClient
	cache  cache.Cache
}

func NewF1Service(client clients.F1DataClient, cache cache.Cache) F1Service {
	return &f1Service{
		client: client,
		cache:  cache,
	}
}

func (service *f1Service) GetScheduleByYear(ctx context.Context, year int) ([]models.RaceWeekend, error) {
	slog.InfoContext(ctx, "Entry: GetScheduleByYear", "year", year)

	cacheKey := fmt.Sprintf("schedule:%d", year)
	if cached, found := service.cache.Get(cacheKey); found {
		slog.InfoContext(ctx, "Cache hit: GetScheduleByYear", "year", year)
		return cached.([]models.RaceWeekend), nil
	}

	if year < 1900 || year > 2050 {
		slog.WarnContext(ctx, "Invalid year requested", "year", year)
		return nil, fmt.Errorf("year outside supported Formula 1 range")
	}

	schedule, err := service.client.GetScheduleByYear(ctx, year)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch schedule", "error", err)
		return nil, fmt.Errorf("failed to fetch schedule from external API: %w", err)
	}

	filteredSchedule := []models.RaceWeekend{}
	for i := range schedule {
		if schedule[i].Round == 0 {
			slog.DebugContext(ctx, "Filtering out non-race event", "name", schedule[i].Name)
			continue
		}
		calculateWeekendBoundaries(&schedule[i])
		formatRaceWeekend(&schedule[i])
		populateStandardCodes(&schedule[i])
		filteredSchedule = append(filteredSchedule, schedule[i])
	}

	service.cache.Set(cacheKey, filteredSchedule, 24*time.Hour)

	slog.InfoContext(ctx, "Exit: GetScheduleByYear", "year", year, "count", len(filteredSchedule))
	return filteredSchedule, nil
}

func (service *f1Service) GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error) {
	slog.InfoContext(ctx, "Entry: GetSessionResults", "year", year, "round", round, "sessionType", sessionType)

	cacheKey := fmt.Sprintf("results:%d:%d:%s", year, round, sessionType)
	if cached, found := service.cache.Get(cacheKey); found {
		slog.InfoContext(ctx, "Cache hit: GetSessionResults", "year", year, "round", round, "sessionType", sessionType)
		return cached.(*models.SessionResults), nil
	}

	results, err := service.client.GetSessionResults(ctx, year, round, sessionType)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch session results", "error", err)
		return nil, fmt.Errorf("failed to fetch results for session %s in round %d (%d): %w", sessionType, round, year, err)
	}

	// Calculate Smart TTL
	ttl := service.calculateSessionTTL(ctx, year, round, sessionType, results != nil && len(results.Results) > 0)

	if results == nil || len(results.Results) == 0 {
		slog.WarnContext(ctx, "No results found", "year", year, "round", round)
		slog.InfoContext(ctx, "Exit: GetSessionResults", "year", year, "round", round, "sessionType", sessionType, "count", 0)

		var finalResults *models.SessionResults
		if results == nil {
			finalResults = &models.SessionResults{
				Year:        year,
				Round:       round,
				SessionType: sessionType,
				Results:     []models.DriverResult{},
			}
		} else {
			finalResults = results
		}

		service.cache.Set(cacheKey, finalResults, ttl)
		return finalResults, nil
	}

	formatSessionResults(results)

	service.cache.Set(cacheKey, results, ttl)

	slog.InfoContext(ctx, "Exit: GetSessionResults", "year", year, "round", round, "sessionType", sessionType, "count", len(results.Results))
	return results, nil
}

func (service *f1Service) calculateSessionTTL(ctx context.Context, year int, round int, sessionType string, hasResults bool) time.Duration {
	if !hasResults {
		return 1 * time.Minute
	}

	// Try to find session start time to determine if it's a current/active session
	schedule, err := service.GetScheduleByYear(ctx, year)
	if err != nil {
		return 10 * time.Minute // Fallback
	}

	var sessionTimeMS int64
	found := false
	for _, weekend := range schedule {
		if weekend.Round == round {
			for _, session := range weekend.Sessions {
				if session.Type == sessionType || formatters.GetSessionCode(session.Type) == sessionType {
					sessionTimeMS = session.TimeUTCMS
					found = true
					break
				}
			}
			break
		}
	}

	if !found {
		return 10 * time.Minute // Fallback
	}

	now := time.Now().UnixMilli()

	// Define "Active Window": from 30 minutes before start to 6 hours after start
	// This covers the session duration and any post-session adjustments.
	if now >= sessionTimeMS-(30*60*1000) && now <= sessionTimeMS+(6*3600*1000) {
		return 1 * time.Minute
	}

	// If it's a historical session (completed > 6 hours ago), cache for a long time
	if now > sessionTimeMS+(6*3600*1000) {
		return 24 * time.Hour
	}

	// Future sessions that aren't in the active window yet
	return 1 * time.Hour
}

func (service *f1Service) GetCircuit(ctx context.Context, year int, round int) (*models.Circuit, error) {
	slog.InfoContext(ctx, "Entry: GetCircuit", "year", year, "round", round)

	cacheKey := fmt.Sprintf("circuit:%d:%d", year, round)
	if cached, found := service.cache.Get(cacheKey); found {
		slog.InfoContext(ctx, "Cache hit: GetCircuit", "year", year, "round", round)
		return cached.(*models.Circuit), nil
	}

	circuit, err := service.client.GetCircuit(ctx, year, round)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch circuit", "error", err)
		return nil, fmt.Errorf("failed to fetch circuit for round %d (%d): %w", round, year, err)
	}

	if circuit != nil {
		circuit.EventDate = formatters.FormatTimestamp(circuit.EventDateMS)
		roundCircuitMetrics(circuit)
		transformLayout(circuit)

		ttl := service.calculateWeekendTTL(ctx, year, round)
		service.cache.Set(cacheKey, circuit, ttl)
	}

	slog.InfoContext(ctx, "Exit: GetCircuit", "year", year, "round", round, "circuit_name", circuit.CircuitName)
	return circuit, nil
}

func (service *f1Service) GetDrivers(ctx context.Context, year int, round int) ([]models.DriverInfo, error) {
	slog.InfoContext(ctx, "Entry: GetDrivers", "year", year, "round", round)

	cacheKey := fmt.Sprintf("drivers:%d:%d", year, round)
	if cached, found := service.cache.Get(cacheKey); found {
		slog.InfoContext(ctx, "Cache hit: GetDrivers", "year", year, "round", round)
		return cached.([]models.DriverInfo), nil
	}

	drivers, err := service.client.GetDrivers(ctx, year, round)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch drivers", "error", err)
		return nil, fmt.Errorf("failed to fetch drivers for round %d (%d): %w", round, year, err)
	}

	if drivers == nil {
		drivers = []models.DriverInfo{}
	}

	ttl := service.calculateWeekendTTL(ctx, year, round)
	service.cache.Set(cacheKey, drivers, ttl)

	slog.InfoContext(ctx, "Exit: GetDrivers", "year", year, "round", round, "count", len(drivers))
	return drivers, nil
}

func (service *f1Service) calculateWeekendTTL(ctx context.Context, year int, round int) time.Duration {
	schedule, err := service.GetScheduleByYear(ctx, year)
	if err != nil {
		return 10 * time.Minute
	}

	var startMS, endMS int64
	found := false
	for _, weekend := range schedule {
		if weekend.Round == round {
			startMS = weekend.StartDateUTCMS
			endMS = weekend.EndDateUTCMS
			found = true
			break
		}
	}

	if !found {
		return 10 * time.Minute
	}

	now := time.Now().UnixMilli()

	// During the weekend (from 24h before start to 12h after end)
	if now >= startMS-(24*3600*1000) && now <= endMS+(12*3600*1000) {
		return 5 * time.Minute
	}

	// Historical
	if now > endMS {
		return 7 * 24 * time.Hour
	}

	// Future
	return 24 * time.Hour
}
